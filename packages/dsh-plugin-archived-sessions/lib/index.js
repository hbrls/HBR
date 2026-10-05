/**
 * Host half of the archived-sessions panel.
 *
 * It owns the join between the two domains so the browser makes a single
 * read-only request: archivedness lives in the Workspace registry, while the
 * session rows (title, activity, running state) live in the Session
 * controller. The route is registered on the official `/api` Fetch extension
 * point, behind the same Host/Origin + browser-auth fence as every other
 * `/api` endpoint.
 *
 * SIMULATION ONLY. Each row also carries the concrete list of things a real
 * delete would remove (session-log files, the projection-cache document, and
 * the workspace accounting slot). Nothing on this half deletes anything: the
 * client prints that plan instead. There is deliberately no destructive
 * endpoint here.
 */

import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';

/** Cordis plugin name. */
export const name = 'plugin-archived-sessions';

/** Required Host services: the /api route registry and both data domains. */
export const inject = ['connection', 'workspaceRegistry', 'sessionController'];

/** The one read-only endpoint this plugin owns. */
export const ROUTE = '/api/archived-sessions';

/** Session ids appear either as a raw uuid or with the `session-` prefix. */
const SESSION_ID_RE = /^(session-)?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/** DSH home directory; the same override every other storage consumer honours. */
function dshHome() {
  return process.env.DSH_HOME || path.join(os.homedir(), '.dsh');
}

/** Root of the per-session JSONL logs: `<root>/<workspace-slug>/<session-id>/`. */
function sessionsRoot() {
  return path.join(dshHome(), 'sessions');
}

/** Per-record projection-cache documents: `<root>/<session-id>.json`. */
function projectionCacheRoot() {
  return path.join(dshHome(), 'storages', 'session_projcache', 'sessions');
}

/** The whole-unit workspace document holding the archive/pin accounting. */
function workspaceStateFile() {
  return path.join(dshHome(), 'storages', 'workspace.json');
}

/**
 * Every spelling one session id may carry across the stores.
 * @param sessionId - the id as the registry spells it.
 * @returns the distinct spellings to look for.
 */
function sessionIdVariants(sessionId) {
  const variants = new Set([sessionId]);
  if (sessionId.startsWith('session-')) variants.add(sessionId.slice('session-'.length));
  else if (SESSION_ID_RE.test(sessionId)) variants.add(`session-${sessionId}`);
  return [...variants];
}

/**
 * Locate every `<root>/<workspace-slug>/<session-id>/` directory. The slug
 * level is scanned rather than re-derived, so the workspace path encoding
 * never has to be reproduced here.
 * @param sessionId - the session to locate.
 * @returns absolute directories that exist on disk.
 */
function findSessionDirs(sessionId) {
  const found = [];
  let entries;
  try {
    entries = fs.readdirSync(sessionsRoot(), { withFileTypes: true });
  } catch {
    return found;
  }
  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    for (const variant of sessionIdVariants(sessionId)) {
      const candidate = path.join(sessionsRoot(), entry.name, variant);
      try {
        if (fs.statSync(candidate).isDirectory()) found.push(candidate);
      } catch {
        /* not this slug */
      }
    }
  }
  return found;
}

/**
 * Every regular file below one directory, depth-first.
 * @param dir - directory to walk.
 * @returns absolute file paths.
 */
function listFilesUnder(dir) {
  const files = [];
  let entries;
  try {
    entries = fs.readdirSync(dir, { withFileTypes: true });
  } catch {
    return files;
  }
  for (const entry of entries) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) files.push(...listFilesUnder(full));
    else files.push(full);
  }
  return files;
}

/**
 * The concrete things a real delete would remove for one session.
 * @param sessionId - the session being deleted.
 * @param archivedSessionIds - the registry-global archive set, for the slot index.
 * @returns ordered targets, each `{ kind, path }`.
 */
function listDeleteTargets(sessionId, archivedSessionIds) {
  const targets = [];
  for (const dir of findSessionDirs(sessionId)) {
    for (const file of listFilesUnder(dir)) targets.push({ kind: 'session-log', path: file });
  }
  for (const variant of sessionIdVariants(sessionId)) {
    const cached = path.join(projectionCacheRoot(), `${variant}.json`);
    try {
      if (fs.statSync(cached).isFile()) targets.push({ kind: 'projection-cache', path: cached });
    } catch {
      /* no cache row for this spelling */
    }
  }
  for (const variant of sessionIdVariants(sessionId)) {
    const slot = archivedSessionIds.indexOf(variant);
    if (slot >= 0) {
      targets.push({
        kind: 'workspace-accounting',
        path: `${workspaceStateFile()} › global.archivedSessionIds[${slot}]`,
      });
    }
  }
  return targets;
}

/**
 * Read every archived session with its workspace and display fields.
 * @param ctx - Host context carrying the Workspace registry and Session controller.
 * @param signal - request lifetime.
 * @returns archived session rows, most recently active first.
 */
async function readArchivedSessions(ctx, signal) {
  const archivedSessionIds = ctx.workspaceRegistry.archivedSessionIds.map(String);

  const workspaceOf = new Map();
  for (const workspace of ctx.workspaceRegistry.list()) {
    for (const sessionId of workspace.sessionIds) workspaceOf.set(String(sessionId), workspace.title);
  }

  const listed = await ctx.sessionController.list({}, signal);
  const rowOf = new Map(listed.items.map((row) => [String(row.sessionId), row]));

  return archivedSessionIds
    .map((sessionId) => {
      const row = rowOf.get(sessionId);
      return {
        sessionId,
        workspaceTitle: workspaceOf.get(sessionId) ?? null,
        title: row?.projections?.values?.title ?? null,
        updatedAt: row?.updatedAt ?? null,
        running: row?.running === true,
        blank: row?.blank === true,
        targets: listDeleteTargets(sessionId, archivedSessionIds),
      };
    })
    .sort((left, right) => (right.updatedAt ?? 0) - (left.updatedAt ?? 0));
}

/**
 * Register the read-only archived-sessions route.
 * @param ctx - Host plugin context.
 */
export function apply(ctx) {
  ctx.effect(() => ctx.connection.fetch.register({
    path: ROUTE,
    methods: ['GET'],
    requestBody: 'buffered',
    fetch: async (request) => {
      const headers = { 'cache-control': 'no-store' };
      try {
        return Response.json({ rows: await readArchivedSessions(ctx, request.signal) }, { headers });
      } catch (error) {
        return Response.json({ error: String(error?.message ?? error) }, { status: 500, headers });
      }
    },
  }), 'plugin-archived-sessions: archived-sessions route');
}
