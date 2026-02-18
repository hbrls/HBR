---
name: accept
description: 检查和验收。包括测试。
---

## Non-Functional Requirements

### Performance

- **NFR-001**: The system shall respond to {{operation}} within {{time}} under {{load conditions}}.
- **NFR-002**: The system shall support {{number}} concurrent {{operations}}.

### Security

- **NFR-010**: The system shall {{security requirement}}.
- **NFR-011**: The system shall NOT {{security prohibition}}.

### Reliability

- **NFR-020**: The system shall maintain {{uptime percentage}} availability.
- **NFR-021**: The system shall recover from {{failure type}} within {{time}}.

### Scalability

- **NFR-030**: The system shall scale to handle {{capacity}} of {{resource}}.

### Usability

- **NFR-040**: The system shall {{usability requirement}}.

### Maintainability

- **NFR-050**: The system shall {{maintainability requirement}}.

## Traceability Matrix

{Maps requirements to acceptance criteria}

| Requirement | Acceptance Criteria | Test Case |
|-------------|---------------------|-----------|
| FR-001 | AC-001 | TC-001 |
| FR-010 | AC-002 | TC-002 |
| NFR-001 | AC-020 | TC-020 |

## Glossary

{Define domain-specific terms}

| Term | Definition |
|------|------------|
| {{term}} | {{definition}} |

## Migration Plan

<!-- If applicable, steps for data migration -->

### Pre-Migration

- [ ] {{step}}

### Migration Steps

1. {{step 1}}
2. {{step 2}}

### Rollback Plan

1. {{rollback step 1}}
2. {{rollback step 2}}

### Post-Migration Verification

- [ ] {{verification step}}

## Testing Strategy

### Unit Tests

- {{what to unit test}}

### Integration Tests

- {{what to integration test}}

### E2E Tests

- {{what to e2e test}}
