# Project: SmartCondo

## Tech Stack

- Astro 7
- TypeScript
- Tailwind CSS 4
- shadcn/ui
- React Hook Form
- Zod

## Architecture

- Place reusable UI primitives in `components/ui`.
- Place feature-specific components in `components/`.
- Shared utilities belong in `lib/`.
- Global types and Zod schemas belong in `types/`.

## Guiding Principles

- Prefer simple, maintainable solutions over clever ones.
- Follow existing project patterns before introducing new ones.
- Reuse existing components and utilities whenever possible.
- Keep components focused on a single responsibility.
- Avoid code duplication.
- Favor composition over inheritance.
- Write self-explanatory code whenever possible.
- Minimize the scope of changes.

## Code Style

- Never use explicit `any`; prefer `unknown` with type guards.
- Use ES Modules (`import` / `export`), never `require()`.
- Use Tailwind CSS only; avoid inline styles and styled-components.
- Add new design tokens to `tailwind.config.ts` before using them.
- Use kebab-case for filenames.
- Use PascalCase for React components.

## Preferred Libraries

Prefer:

- Astro Actions
- React Hook Form
- Zod
- shadcn/ui

Avoid introducing new libraries unless there is a clear benefit.

## Workflow

- Always run `npm run type-check && npm run lint` after a series of changes.
- Run targeted tests instead of the full suite whenever possible:
  `npm run test -- FileName`
- Branch names:
  - `feat/...`
  - `fix/...`
  - `chore/...`
- Write commit messages in English using the imperative mood.
  Example: `add OAuth callback handler`

## Development

Start the development server in background mode:

```sh
astro dev --background
```

Manage the background server with:

```sh
astro dev status
astro dev logs
astro dev stop
```

## Restrictions

- Do not introduce new dependencies without justification.
- Do not change the project structure.
- Preserve backward compatibility.
- Keep changes minimal.
- Do not refactor unrelated code while implementing a feature.
- Do not ignore TypeScript or lint errors.

## Before Finishing

Before completing a task:

- Run type checking.
- Run linting.
- Update types if necessary.
- Verify imports.
- Remove unused code and imports.
- Explain any important architectural decisions.
- Confirm the implementation follows the project conventions.

## Documentation

Official documentation:

https://docs.astro.build

Consult these guides before working on related tasks:

- Routing, pages and middleware
- Astro components
- Framework components
- Content collections
- Styling and Tailwind CSS
- Internationalization
