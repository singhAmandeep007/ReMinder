import { Code, Server, TestTube, Settings } from "lucide-react";

import { TTechStackData } from "./types";

export const techStack: TTechStackData = {
  frontend: {
    name: "Frontend",
    color: "#22c55e",
    icon: <Code className="h-6 w-6" />,
    description: "User interface and client-side logic",
    technologies: [
      {
        id: "react",
        name: "React",
        description: "UI library",
        icon: "https://singhamandeep007.github.io/logos/react.svg",
      },
      {
        id: "typescript",
        name: "TypeScript",
        description: "Type-safe JavaScript",
        icon: "https://singhamandeep007.github.io/logos/typescript.svg",
      },
      {
        id: "redux",
        name: "Redux Toolkit",
        description: "State management",
        icon: "https://singhamandeep007.github.io/logos/redux.svg",
      },
      { id: "shadcn", name: "Shadcn UI", description: "Component library" },
      {
        id: "tailwind",
        name: "Tailwind CSS",
        description: "Utility-first CSS",
        icon: "https://singhamandeep007.github.io/logos/tailwind.svg",
      },
      { id: "router", name: "React Router", description: "Routing" },
      { id: "i18n", name: "i18n", description: "Localization" },
    ],
  },
  backend: {
    name: "Backend",
    color: "#3b82f6",
    icon: <Server className="h-6 w-6" />,
    description: "Server-side logic and data processing",
    technologies: [
      {
        id: "nestjs",
        name: "NestJS",
        description: "TypeScript framework",
        icon: "https://docs.nestjs.com/assets/logo-small-gradient.svg",
      },
      {
        id: "golang",
        name: "Golang",
        description: "Go language",
        icon: "https://singhamandeep007.github.io/logos/golang.svg",
      },
      { id: "postgres", name: "PostgreSQL", description: "Relational database" },
      {
        id: "firestore",
        name: "Firestore",
        description: "NoSQL database",
        icon: "https://www.gstatic.com/devrel-devsite/prod/v6dc4611c4232bd02b2b914c4948f523846f90835f230654af18f87f75fe9f73c/firebase/images/lockup.svg",
      },
      {
        id: "mongodb",
        name: "MongoDB",
        description: "NoSQL database",
        icon: "https://singhamandeep007.github.io/logos/mongodb.svg",
      },
      { id: "sqlite", name: "SQLite", description: "Embedded database" },
      { id: "gin", name: "Gin", description: "Go web framework" },
    ],
  },
  testing: {
    name: "Testing",
    color: "#f59e0b",
    icon: <TestTube className="h-6 w-6" />,
    description: "Quality assurance and validation",
    technologies: [
      {
        id: "jest",
        name: "Jest",
        description: "Testing framework",
        icon: "https://singhamandeep007.github.io/logos/jest.svg",
      },
      {
        id: "rtl",
        name: "React Testing Library",
        description: "Component testing",
        icon: "https://testing-library.com/img/octopus-64x64.png",
      },
      {
        id: "cypress",
        name: "Cypress",
        description: "E2E testing",
        icon: "https://singhamandeep007.github.io/logos/cypress.svg",
      },
      {
        id: "playwright",
        name: "Playwright",
        description: "Browser automation",
        icon: "https://singhamandeep007.github.io/logos/playwright.svg",
      },
      {
        id: "msw",
        name: "MSW",
        description: "API mocking",
        icon: "https://singhamandeep007.github.io/logos/msw.svg",
      },
      {
        id: "mirage",
        name: "Mirage.js",
        description: "Mock server",
        icon: "https://avatars.githubusercontent.com/u/47899903?s=48&v=4",
      },
      {
        id: "istanbul",
        name: "Istanbul",
        description: "Code coverage",
        icon: "https://singhamandeep007.github.io/logos/css.svg",
      },
    ],
  },
  devtools: {
    name: "Dev Tools",
    color: "#8b5cf6",
    icon: <Settings className="h-6 w-6" />,
    description: "Development workflow and tooling",
    technologies: [
      {
        id: "storybook",
        name: "Storybook",
        description: "Component development",
        icon: "https://singhamandeep007.github.io/logos/storybook.svg",
      },
      { id: "eslint", name: "ESLint", description: "Code linting" },
      { id: "prettier", name: "Prettier", description: "Code formatting" },
      { id: "husky", name: "Husky", description: "Git hooks" },
      { id: "commitizen", name: "Commitizen", description: "Commit conventions" },
      {
        id: "typescript-tools",
        name: "TypeScript",
        description: "Static typing",
        icon: "https://singhamandeep007.github.io/logos/typescript.svg",
      },
      {
        id: "webpack",
        name: "Webpack",
        description: "Module bundler",
        icon: "https://raw.githubusercontent.com/webpack/media/90b54d02fa1cfc8aa864a8322202f74ac000f5d2/logo/icon.svg",
      },
    ],
  },
};
