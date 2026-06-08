---
name: template-design
metadata:
    version: 0.0.2
---

# Design Principles

## Design & Architecture

### Languages

| Purpose | Language | Version | Notes |
|---------|----------|---------|-------|
| Backend | {{language}} | {{version}} | |
| Frontend | {{language}} | {{version}} | |
| Scripts | {{language}} | {{version}} | |

### Frameworks

| Purpose | Framework | Version | Notes |
|---------|-----------|---------|-------|
| Web Framework | {{framework}} | {{version}} | |
| ORM | {{framework}} | {{version}} | |
| Testing | {{framework}} | {{version}} | |

### Infrastructure

| Component | Design Choice | Provider | Notes |
|-----------|------------|----------|-------|
| Database | {{design_choice}} | {{provider}} | |
| Cache | {{design_choice}} | {{provider}} | |
| Message Queue | {{design_choice}} | {{provider}} | |
| Hosting | {{design_choice}} | {{provider}} | |

## Performance Targets

| Metric | Target | Current | Priority |
|--------|--------|---------|----------|
| API Response Time (p50) | {{ms}} | {{ms}} | High |
| API Response Time (p99) | {{ms}} | {{ms}} | High |
| Page Load Time | {{s}} | {{s}} | High |
| Memory Usage | {{MB}} | {{MB}} | Medium |
| CPU Utilization | {{%}} | {{%}} | Medium |

## Scalability Requirements

- **Concurrent Users**: {{number}}
- **Requests per Second**: {{number}}
- **Data Volume**: {{size}}
- **Growth Rate**: {{percentage}} per {{period}}

## Security Requirements

### Authentication

- {{authentication method}}
- {{session management}}

### Authorization

- {{authorization model}}
- {{role definitions}}

### Data Protection

- **Encryption at Rest**: {{requirement}}
- **Encryption in Transit**: {{requirement}}
- **PII Handling**: {{requirement}}

### Compliance

- {{compliance requirement 1}}
- {{compliance requirement 2}}

## Dependencies

### External Services

| Service | Purpose | Criticality | Fallback |
|---------|---------|-------------|----------|
| {{service}} | {{purpose}} | High/Medium/Low | {{fallback}} |

### Internal Services

| Service | Purpose | Owner | API Docs |
|---------|---------|-------|----------|
| {{service}} | {{purpose}} | {{team}} | {{link}} |

## Development Environment

### Required Tools

- {{tool 1}} (version {{version}})
- {{tool 2}} (version {{version}})

### Setup

<!-- // 遵循团队统一的开发环境配置，如有特殊安装或启动步骤，在此说明。 -->

## Design Debt

### Known Issues

| Issue | Impact | Priority | Plan |
|-------|--------|----------|------|
| {{issue}} | {{impact}} | High/Medium/Low | {{resolution plan}} |

### Design Debt Budget

- {{percentage}}% of sprint capacity allocated to tech debt
- Major refactoring requires spec

## Constraints

### Must Use

<!-- Designs that must be used -->

- {{design}}: {{reason}}

### Must Not Use

<!-- Designs that are prohibited -->

- {{design}}: {{reason}}

### Prefer

<!-- Designs to prefer when applicable -->

- {{design}}: {{reason}}

## Architecture Decisions Records (ADRs)

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| ADR-001 | {{title}} | Accepted | {{date}} |
| ADR-002 | {{title}} | Accepted | {{date}} |

## Changelog

<!-- // 这是一个 Living Document，如无必要，无需维护变更历史。 -->
