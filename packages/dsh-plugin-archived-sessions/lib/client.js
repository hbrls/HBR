window.__ModuleLoader__.load({
  id: 'dsh-plugin-archived-sessions',
  factory: (require) => {
    const React = require('react');
    const {
      Button, IconArchiveOutlineRegular, IconTrashOutlineRegular, IconWarningOutlineRegular,
      MarkdownText, Modal, relativeTime,
    } = require('@deepseek-ai/dsh-client-ui-primitives');

    const h = React.createElement;
    /** Locale namespace owned by this plugin. */
    const NS = 'plugin-archived-sessions';
    /** The single read-only route owned by this plugin's Host half. */
    const ROUTE = '/api/archived-sessions';
    /** How long one simulated deletion line takes to appear. */
    const STEP_MS = 220;
    /** Dialog class names, passed through the official Modal `className` props. */
    const DLG = 'dsh-archived-sessions-delete-dialog';
    const DLG_CONTENT = `${DLG}-content`;

    /**
     * The official confirmation-dialog geometry, copied from the primitives'
     * own `RiskConfirmation` stylesheet — the card caps its own growth and
     * clips, the content column scrolls, and `100dvh` keeps it under the
     * window — with one deliberate difference: a wider card, because this
     * dialog lists absolute file paths.
     */
    const DIALOG_CSS = `
.${DLG} { width: min(720px, 100%); max-height: 100%; overflow: hidden; }
.${DLG_CONTENT} { min-height: 0; overflow-y: auto; overscroll-behavior: contain; }
@supports (height: 100dvh) { .${DLG} { max-height: calc(100dvh - 48px); } }
.${DLG}-warning { display: flex; align-items: flex-start; gap: 10px; color: var(--dsw-alias-label-secondary); font-size: 14px; line-height: 22px; }
.${DLG}-warning p { margin: 0; }
.${DLG}-warning-icon { flex: none; margin-top: 2px; color: var(--dsw-alias-state-error-primary); }
.${DLG}-list { display: flex; flex-direction: column; gap: 8px; margin-top: 20px; }
.${DLG}-line { display: flex; align-items: baseline; gap: 10px; }
.${DLG}-num { flex: none; min-width: 22px; text-align: right; color: var(--dsw-alias-label-tertiary); font-size: 13px; line-height: 20px; }
.${DLG}-path { flex: 1; min-width: 0; font-family: var(--ds-font-family-code); font-size: 13px; line-height: 20px; color: var(--dsw-alias-label-primary); overflow-wrap: anywhere; }
.${DLG}-state { flex: none; color: var(--dsw-alias-state-success-primary); font-size: 13px; line-height: 20px; white-space: nowrap; }
.${DLG}-restart { margin-top: 20px; padding-top: 16px; border-top: 0.5px solid var(--dsw-alias-border-l1); color: var(--dsw-alias-state-warn-primary); font-size: 14px; line-height: 22px; font-weight: 500; }
.${DLG}-note { margin-top: 12px; color: var(--dsw-alias-label-tertiary); font-size: 12px; line-height: 18px; }
`;

    /**
     * Inject the dialog stylesheet once. It is appended to `document.head`
     * after the shell's own stylesheet has loaded, which is what lets the
     * width rule below win over the official `.dialog` rule: both are single
     * class selectors, so document order decides. This runs at plugin apply
     * time, after boot.
     */
    let stylesInstalled = false;
    function installStyles() {
      if (stylesInstalled || typeof document === 'undefined') return;
      stylesInstalled = true;
      const tag = document.createElement('style');
      tag.setAttribute('data-dsh-plugin', NS);
      tag.textContent = DIALOG_CSS;
      document.head.appendChild(tag);
    }

    const zh = {
      'col.workspace': '工作区',
      'col.session': '会话标题',
      'col.state': '状态',
      'col.updated': '最近活动',
      'col.id': '会话 ID',
      'col.action': '操作',
      'state.running': '运行中',
      'state.blank': '空白',
      'state.archived': '已归档',
      'session.untitled': '未命名',
      'workspace.none': '未归属工作区',
      'loading': '读取中…',
      'empty': '没有已归档会话',
      'error': '读取失败',
      'time.now': '刚刚',
      'time.minutes': '{n}分钟',
      'time.hours': '{n}小时',
      'time.days': '{n}天',
      'time.months': '{n}个月',
      'time.years': '{n}年',
      'action.delete': '删除',
      'dialog.title': '删除会话',
      'dialog.intro': '我们将删除以下几个文件',
      'dialog.none': '没有找到可删除的文件',
      'dialog.running': '正在删除…',
      'dialog.deleted': '已删除',
      'dialog.restart': '请重启',
      'dialog.cancel': '取消',
      'dialog.confirm': '删除',
      'dialog.close': '关闭',
      'dialog.simulated': '模拟删除：不会真正修改任何文件。',
    };

    const en = {
      'col.workspace': 'Workspace',
      'col.session': 'Session',
      'col.state': 'State',
      'col.updated': 'Last active',
      'col.id': 'Session ID',
      'col.action': 'Actions',
      'state.running': 'Running',
      'state.blank': 'Blank',
      'state.archived': 'Archived',
      'session.untitled': 'Untitled',
      'workspace.none': 'No workspace',
      'loading': 'Loading…',
      'empty': 'No archived sessions',
      'error': 'Read failed',
      'time.now': 'now',
      'time.minutes': '{n}min',
      'time.hours': '{n}h',
      'time.days': '{n}d',
      'time.months': '{n}mo',
      'time.years': '{n}y',
      'action.delete': 'Delete',
      'dialog.title': 'Delete session',
      'dialog.intro': 'The following files will be deleted',
      'dialog.none': 'Nothing to delete was found',
      'dialog.running': 'Deleting…',
      'dialog.deleted': 'deleted',
      'dialog.restart': 'Please restart',
      'dialog.cancel': 'Cancel',
      'dialog.confirm': 'Delete',
      'dialog.close': 'Close',
      'dialog.simulated': 'Simulated deletion: no file is actually changed.',
    };

    /**
     * The sidebar panel glyph. The sidebar renders every `sidebar.panellist`
     * entry as `{ size, active }` at 16px (rail) or 18px (wide) and owns the
     * active styling itself, so the entry only forwards the official icon at
     * the requested size — the shape every shipped panel icon uses.
     */
    const Icon = ({ size }) => h(IconArchiveOutlineRegular, { size });

    /**
     * Compact relative time through the official primitives bucketing and the
     * official locale copy (same shape as dsh-client-ui-workspace's timeLabel).
     */
    function timeLabel(updatedAt, now, t) {
      if (typeof updatedAt !== 'number') return '—';
      const { unit, n } = relativeTime(updatedAt, now);
      return unit === 'now' ? t('time.now') : t(`time.${unit}`, { n });
    }

    const wrapStyle = { flex: 1, height: '100%', overflow: 'auto', padding: '24px' };
    const tableStyle = { borderCollapse: 'collapse', width: '100%', fontSize: 13, lineHeight: '20px' };
    const thStyle = {
      textAlign: 'left',
      padding: '6px 10px',
      borderBottom: '1px solid var(--dsw-alias-border-l2, rgba(128,128,128,.35))',
      color: 'var(--dsw-alias-label-secondary, #8a8a8e)',
      fontWeight: 500,
      whiteSpace: 'nowrap',
    };
    const tdStyle = {
      padding: '6px 10px',
      borderBottom: '1px solid var(--dsw-alias-border-l1, rgba(128,128,128,.18))',
      color: 'var(--dsw-alias-label-primary, inherit)',
      verticalAlign: 'middle',
    };
    const idStyle = {
      fontFamily: 'var(--ds-font-family-code)',
      fontSize: 12,
      color: 'var(--dsw-alias-label-secondary)',
    };

    /**
     * One read-only request. The Host half owns the join between the Workspace
     * registry (which session is archived, and its workspace) and the Session
     * controller (title, activity, running state), and also resolves the files
     * a real delete would remove.
     */
    async function readArchivedSessions(signal, t) {
      const response = await fetch(ROUTE, {
        method: 'GET',
        headers: { accept: 'application/json' },
        signal,
      });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const body = await response.json();
      if (body.error !== undefined) throw new Error(String(body.error));
      return body.rows.map((row) => ({
        ...row,
        workspaceTitle: row.workspaceTitle ?? t('workspace.none'),
        targets: Array.isArray(row.targets) ? row.targets : [],
      }));
    }

    /** The archived rows as a real table, so each row can carry an action. */
    function SessionTable({ rows, t, now, onDelete }) {
      return h('table', { style: tableStyle },
        h('thead', null, h('tr', null,
          h('th', { style: thStyle }, t('col.workspace')),
          h('th', { style: thStyle }, t('col.session')),
          h('th', { style: thStyle }, t('col.state')),
          h('th', { style: thStyle }, t('col.updated')),
          h('th', { style: thStyle }, t('col.id')),
          h('th', { style: thStyle }, t('col.action')))),
        h('tbody', null, rows.map((row) => {
          const state = row.running ? t('state.running') : row.blank ? t('state.blank') : t('state.archived');
          const empty = row.targets.length === 0;
          return h('tr', { key: row.sessionId },
            h('td', { style: tdStyle }, row.workspaceTitle),
            h('td', { style: tdStyle }, row.title ?? t('session.untitled')),
            h('td', { style: tdStyle }, state),
            h('td', { style: tdStyle }, timeLabel(row.updatedAt, now, t)),
            h('td', { style: tdStyle }, h('code', { style: idStyle }, String(row.sessionId).slice(-8))),
            h('td', { style: tdStyle }, h(Button, {
              variant: 'outline',
              size: 'sm',
              icon: h(IconTrashOutlineRegular, { size: 14 }),
              disabled: empty,
              onClick: () => onDelete(row),
            }, t('action.delete'))));
        })));
    }

    /**
     * The delete dialog. It first lists what a real delete would remove, then
     * walks that list printing one `N. <path> 已删除` line per step and ends
     * with `请重启`. Nothing is sent to the Host: this is a simulation.
     */
    function DeleteDialog({ row, t, onClose }) {
      const targets = row.targets;
      const [phase, setPhase] = React.useState('confirm');
      const [done, setDone] = React.useState(0);

      React.useEffect(() => {
        if (phase !== 'running') return undefined;
        if (done >= targets.length) {
          setPhase('done');
          return undefined;
        }
        const timer = setTimeout(() => setDone((n) => n + 1), STEP_MS);
        return () => clearTimeout(timer);
      }, [phase, done, targets.length]);

      const start = React.useCallback(() => {
        setDone(0);
        setPhase('running');
      }, []);

      const requestClose = React.useCallback(() => {
        if (phase === 'running') return;
        onClose();
      }, [phase, onClose]);

      const finished = phase === 'done';
      // Only the confirm and running phases carry a notice; the finished phase
      // states the restart instead, so `请重启` is never printed twice.
      const notice = phase === 'confirm'
        ? t('dialog.intro')
        : phase === 'running'
          ? t('dialog.running')
          : null;

      const footer = phase === 'confirm'
        ? [
          h(Button, { key: 'cancel', variant: 'outline', onClick: requestClose }, t('dialog.cancel')),
          h(Button, {
            key: 'confirm',
            variant: 'primary',
            disabled: targets.length === 0,
            onClick: start,
          }, t('dialog.confirm')),
        ]
        : [
          h(Button, {
            key: 'close',
            variant: 'outline',
            disabled: phase === 'running',
            onClick: requestClose,
          }, t('dialog.close')),
        ];

      const body = [];
      if (notice !== null) {
        body.push(h('div', { key: 'warning', className: `${DLG}-warning` },
          h(IconWarningOutlineRegular, { size: 18, className: `${DLG}-warning-icon` }),
          h('p', null, notice)));
      }
      body.push(h('div', { key: 'list', className: `${DLG}-list` },
        targets.length === 0
          ? h('div', { className: `${DLG}-note` }, t('dialog.none'))
          : targets.map((target, index) => h('div', { key: index, className: `${DLG}-line` },
            h('span', { className: `${DLG}-num` }, `${index + 1}.`),
            h('code', { className: `${DLG}-path` }, target.path),
            phase !== 'confirm' && index < done
              ? h('span', { className: `${DLG}-state` }, t('dialog.deleted'))
              : null))));
      if (finished) {
        body.push(h('div', { key: 'restart', className: `${DLG}-restart` }, t('dialog.restart')));
      }
      body.push(h('div', { key: 'note', className: `${DLG}-note` }, t('dialog.simulated')));

      return h(Modal, {
        open: true,
        onClose: requestClose,
        title: t('dialog.title'),
        closeLabel: t('dialog.cancel'),
        className: DLG,
        contentClassName: DLG_CONTENT,
        footer,
      }, body);
    }

    /** The main-area page: the panel title, the table, and the dialog. */
    function ArchivedSessionsPage({ ctx }) {
      const t = React.useMemo(() => ctx.locale.bind(NS), [ctx]);
      const [state, setState] = React.useState({ status: 'loading' });
      const [pending, setPending] = React.useState(null);

      React.useEffect(() => {
        const lifetime = new AbortController();
        readArchivedSessions(lifetime.signal, t).then(
          (rows) => {
            if (!lifetime.signal.aborted) setState({ status: 'ready', rows });
          },
          (error) => {
            if (!lifetime.signal.aborted) {
              setState({ status: 'error', message: String(error?.message ?? error) });
            }
          },
        );
        return () => {
          lifetime.abort();
        };
      }, [t]);

      if (state.status === 'loading') {
        return h('div', { style: wrapStyle }, h(MarkdownText, { text: `# Archived Sessions\n\n${t('loading')}` }));
      }
      if (state.status === 'error') {
        return h('div', { style: wrapStyle }, h(MarkdownText, { text: `# Archived Sessions\n\n${t('error')}: ${state.message}` }));
      }
      if (state.rows.length === 0) {
        return h('div', { style: wrapStyle }, h(MarkdownText, { text: `# Archived Sessions\n\n${t('empty')}` }));
      }

      return h('div', { style: wrapStyle },
        h(MarkdownText, { text: '# Archived Sessions' }),
        h(SessionTable, { rows: state.rows, t, now: Date.now(), onDelete: setPending }),
        pending ? h(DeleteDialog, { row: pending, t, onClose: () => setPending(null) }) : null);
    }

    return {
      inject: ['slots', 'locale'],
      apply(ctx) {
        installStyles();
        ctx.effect(() => ctx.locale.register(NS, { zh, en }), 'plugin-archived-sessions: dictionaries');

        ctx.slots.inject('sidebar.panellist', () => ctx.slots.register({
          name: 'sidebar.panellist', id: 'plugin-archived-sessions', order: 1, label: 'Archived Sessions',
        }, Icon));

        ctx.slots.inject('main', () => ctx.slots.register({
          name: 'main', key: 'plugin-archived-sessions',
        }, () => h(ArchivedSessionsPage, { ctx })));

      },
    };
  },
});
