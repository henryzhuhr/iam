---
name: requirements-analysis
description: Analyze a product idea, feature request, business requirement, or project brief into a structured development-ready requirements breakdown. Use when the user wants Codex to turn a rough idea into clear modules, scope, user roles, workflows, functional requirements, non-functional requirements, risks, dependencies, open questions, priorities, or phased implementation tasks.
---

# Requirements Analysis

## Overview

Turn vague or incomplete需求 into a consistent analysis output that is immediately useful for product planning and engineering delivery. Prefer Chinese output when the user writes in Chinese.

## Workflow

1. Read the user's input and restate the target product, feature, or business goal in one short paragraph.
2. Infer the scenario, target users, and expected business value from the available context.
3. Separate confirmed facts from assumptions. If key information is missing, do not block by default; continue with reasonable assumptions and mark them clearly.
4. Break the requirement into modules or capability areas that are suitable for implementation planning.
5. Expand each module into concrete requirements, constraints, and acceptance expectations.
6. Identify cross-cutting concerns such as permissions, data model, integration, security, performance, observability, and rollout risk.
7. End with a prioritized implementation view so the result can feed PRD, task breakdown, API design, or development scheduling.

## Output Rules

- Use the standard structure from `references/output-template.md`.
- Use the dimension checklist from `references/analysis-dimensions.md` to avoid missing important areas.
- Keep the response structured and actionable rather than essay-like.
- Prefer explicit labels such as `已确认`, `假设`, `待确认`, `建议优先级`.
- If the user's requirement is still fuzzy, produce a usable first-pass analysis plus a short list of clarification questions.
- If the request is narrow, compress the output but still cover scope, requirements, constraints, and acceptance criteria.
- If the request is broad, group results by phase or module to keep the output readable.

## Default Behavior

### For incomplete input

- Infer a reasonable product boundary.
- Mark uncertain items as assumptions.
- Highlight the minimum questions that would materially change implementation.

### For development-oriented requests

- Emphasize modules, APIs, data objects, permissions, edge cases, and acceptance criteria.
- Add a suggested implementation order with MVP first.

### For planning-oriented requests

- Emphasize business goals, user roles, scenarios, scope boundaries, and priority.
- Keep technical details high-level unless the user asks for architecture depth.

## Quality Bar

- Distinguish user goals from implementation ideas.
- Avoid mixing mandatory requirements with optional enhancements.
- Convert ambiguous wording into verifiable statements whenever possible.
- Surface hidden complexity early: state transitions, authorization, data consistency, third-party dependencies, migration impact, auditability, and failure handling.
- Produce output that another engineer can directly use as the basis for PRD, issue breakdown, or milestone planning.

## References

- Read `references/output-template.md` for the default response format.
- Read `references/analysis-dimensions.md` when the scope is medium or large, or when the user input is ambiguous and you need a fuller checklist.
