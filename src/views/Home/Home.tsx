import { FC } from "react";

import { Code2, Cpu, Github, Layers, Pickaxe, Settings2, Shield, Sparkles, TestTube2 } from "lucide-react";

import { Button, RouteLink, Card, CardContent } from "components";

import { HOME_ROUTE_BY_PATH } from "app/Router";

import { TechStackVisualization } from "./TechStackVisualization";
import { DependencyVisualization } from "./DependencyVisualization";
import { BuildProcessVisualization } from "./BuildProcessVisualization";

const technicalFeatures = [
  {
    icon: (
      <Layers
        size={32}
        className="mb-4 text-primary"
      />
    ),
    title: "Multi Backend Implementation",
    description:
      "Showcasing two robust backend implementations: a NestJS TypeScript server and a Go Gin server, both following clean architecture principles and RESTful API design.",
  },
  {
    icon: (
      <Shield
        size={32}
        className="mb-4 text-primary"
      />
    ),
    title: "Type-Safe Development",
    description:
      "Leveraging TypeScript for robust type safety across the entire stack, ensuring reliability and better developer experience throughout the codebase.",
  },
  {
    icon: (
      <Code2
        size={32}
        className="mb-4 text-primary"
      />
    ),
    title: "Modern Frontend Architecture",
    description:
      "Built with React 18, Redux Toolkit for state management, and Tailwind CSS for styling. Features comprehensive testing with Jest, React Testing Library, and Cypress.",
  },
  {
    icon: (
      <TestTube2
        size={32}
        className="mb-4 text-primary"
      />
    ),
    title: "Quality Assurance",
    description:
      "Comprehensive testing strategy including unit tests (Jest), integration tests (React Testing Library), E2E tests (Cypress), and Storybook for component development.",
  },
];

const developmentSteps = [
  {
    icon: (
      <Layers
        size={24}
        className="mr-3 text-primary"
      />
    ),
    title: "Clean Architecture",
    detail:
      "Implemented domain-driven design with clear separation of concerns across both backend implementations and the frontend application.",
  },
  {
    icon: (
      <Settings2
        size={24}
        className="mr-3 text-primary"
      />
    ),
    title: "Modern Tooling",
    detail:
      "Utilizing cutting-edge tools: TypeScript, ESLint, Prettier, Husky for git hooks, and comprehensive CI/CD setup.",
  },
  {
    icon: (
      <TestTube2
        size={24}
        className="mr-3 text-primary"
      />
    ),
    title: "Quality First",
    detail: "Built with a test-first approach, featuring comprehensive test coverage and automated quality checks.",
  },
];

// Technical achievements
const technicalAchievements = [
  {
    quote:
      "The project demonstrates excellent use of TypeScript for type safety and maintainability. The architecture is clean and well-organized, with clear separation of concerns.",
    author: "Code Review",
    icon: <Layers className="h-12 w-12 text-primary" />,
  },
  {
    quote:
      "Impressive test coverage and implementation of modern React patterns. The state management solution using Redux Toolkit is particularly well-designed.",
    author: "Technical Assessment",
    icon: <Code2 className="h-12 w-12 text-primary" />,
  },
  {
    quote:
      "A great example of how to build a maintainable and scalable full-stack application with proper separation of concerns and modern development practices.",
    author: "Architecture Review",
    icon: <TestTube2 className="h-12 w-12 text-primary" />,
  },
];

export type THomeProps = Record<string, never>;

export const Home: FC<THomeProps> = () => {
  return (
    <div data-testid="home-page">
      <header className="py-12 text-center sm:py-20">
        <div className="container mx-auto px-4 sm:px-6">
          <h1 className="mb-6 bg-gradient-to-r from-primary via-primary/80 to-primary/60 bg-clip-text text-4xl font-extrabold text-transparent sm:text-5xl md:text-7xl">
            Full-Stack Reminder Application
          </h1>
          <p className="mx-auto mb-10 max-w-3xl text-lg text-muted-foreground sm:text-xl md:text-2xl">
            A production-grade reminder application showcasing modern web development practices, featuring multi backend
            implementations in Go, TypeScript and Python, with a robust React frontend.
          </p>
          <div className="flex flex-col items-center justify-center gap-4 sm:flex-row sm:space-x-4">
            <Button
              asChild
              size="lg"
              className="w-full transform px-5 py-4 text-lg font-bold shadow-xl transition duration-300 ease-in-out hover:scale-105 sm:w-auto"
            >
              <RouteLink to={HOME_ROUTE_BY_PATH.reminders}>
                <div className="flex items-center gap-1">
                  <Sparkles
                    className="-mt-1 mr-2 inline-block"
                    size={20}
                  />
                  Explore Demo
                </div>
              </RouteLink>
            </Button>
            <Button
              asChild
              variant="outline"
              size="lg"
              className="w-full transform border-2 px-5 py-4 text-lg font-semibold transition duration-300 ease-in-out hover:scale-105 sm:w-auto"
            >
              <a href="#build">
                <div className="flex items-center gap-1">
                  <Pickaxe
                    className="-mt-1 mr-2 inline-block"
                    size={20}
                  />
                  Check Build Process
                </div>
              </a>
            </Button>
          </div>
          <div className="mx-auto mt-16 h-[600px] w-full max-w-5xl rounded-xl border border-border bg-card p-4 md:mt-24">
            <TechStackVisualization />
          </div>
        </div>
      </header>

      <section className="border-t border-border bg-muted py-12 sm:py-16 md:py-24">
        <div className="container mx-auto px-4 text-center sm:px-6">
          <h2 className="mb-6 text-2xl font-bold sm:text-3xl md:text-4xl">
            The Challenge: Building a Scalable Reminder System
          </h2>
          <p className="mx-auto mb-8 max-w-2xl text-base text-muted-foreground sm:text-lg">
            Creating a reminder application that's both powerful and maintainable requires careful consideration of
            architecture, state management, and user experience. This project demonstrates how to build such a system
            using modern web technologies, with a focus on scalability and maintainability.
          </p>
          <p className="bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-lg font-semibold text-transparent sm:text-xl">
            Enter ReMinder: A Case Study in Modern Full-Stack Development
          </p>
        </div>
      </section>

      {/* Frontend Dependencies Visualization */}
      <section className="py-12 sm:py-16 md:py-24">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-8 text-center text-3xl font-bold sm:text-4xl md:text-5xl">Frontend Dependencies</h2>
          <p className="mx-auto mb-12 max-w-2xl text-center text-muted-foreground">
            An interactive visualization of the project's frontend dependencies, categorized by their purpose in the
            application architecture.
          </p>
          <div className="mx-auto max-w-5xl rounded-xl border border-border bg-card p-6 shadow-lg">
            <DependencyVisualization />
          </div>
        </div>
      </section>

      <section
        className="border-t border-border bg-muted py-12 sm:py-16 md:py-24"
        id="build"
      >
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-8 text-center text-3xl font-bold sm:text-4xl md:text-5xl">Build Process Visualization</h2>
          <p className="mx-auto mb-12 max-w-2xl text-center text-muted-foreground">
            Watch how the application is built from source code to deployable package, illustrating the transformation
            process.
          </p>
          <div className="mx-auto max-w-5xl rounded-xl border border-border bg-card p-6 shadow-lg">
            <BuildProcessVisualization />
          </div>
        </div>
      </section>

      {/* Technical Features Section */}
      <section
        id="features"
        className="py-12 sm:py-16 md:py-24"
      >
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-12 text-center text-3xl font-bold sm:mb-16 sm:text-4xl md:text-5xl">
            Technical Highlights
          </h2>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {technicalFeatures.map((feature, index) => (
              <Card
                key={index}
                className="transform transition-all duration-300 hover:scale-105 hover:shadow-primary/20"
              >
                <CardContent className="p-6">
                  <div className="flex justify-center md:justify-start">{feature.icon}</div>
                  <h3 className="mb-3 text-xl font-semibold">{feature.title}</h3>
                  <p className="text-sm text-muted-foreground">{feature.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Development Approach Section */}
      <section className="border-t border-border bg-muted py-12 sm:py-16 md:py-24">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-12 text-center text-3xl font-bold sm:mb-16 sm:text-4xl md:text-5xl">
            Development Approach
          </h2>
          <div className="relative mx-auto max-w-4xl">
            <div className="absolute left-0 top-1/2 hidden h-0.5 w-full -translate-y-1/2 bg-border md:block"></div>
            <div className="relative grid gap-6 md:grid-cols-3">
              {developmentSteps.map((step, index) => (
                <Card
                  key={index}
                  className="z-10 transform transition-transform duration-300 hover:z-20 hover:scale-105"
                >
                  <CardContent className="p-6 text-center">
                    <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-muted ring-4 ring-border">
                      {step.icon}
                    </div>
                    <h3 className="mb-2 text-lg font-semibold">
                      {index + 1}. {step.title}
                    </h3>
                    <p className="text-sm text-muted-foreground">{step.detail}</p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Technical Achievements Section */}
      <section className="py-12 sm:py-16 md:py-24">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-12 text-center text-3xl font-bold sm:mb-16 sm:text-4xl md:text-5xl">
            Technical Achievements
          </h2>
          <div className="grid gap-6 sm:grid-cols-2 md:grid-cols-3">
            {technicalAchievements.map((achievement, index) => (
              <Card
                key={index}
                className="flex transform flex-col items-center transition-transform duration-300 hover:-translate-y-2"
              >
                <CardContent className="flex flex-col items-center p-6 text-center">
                  <div className="mb-6">{achievement.icon}</div>
                  <blockquote className="mb-6 flex-grow text-base italic text-muted-foreground sm:text-lg">
                    "{achievement.quote}"
                  </blockquote>
                  <p className="font-semibold">{achievement.author}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Tech Stack Section */}
      <section
        id="tech-stack"
        className="border-t border-border bg-muted py-12 sm:py-16 md:py-24"
      >
        <div className="container mx-auto max-w-3xl px-4 text-center sm:px-6">
          <Cpu
            size={48}
            className="mx-auto mb-6 text-primary"
          />
          <h2 className="mb-6 text-2xl font-bold sm:text-3xl md:text-4xl">Technology Stack</h2>
          <p className="mb-4 text-base text-muted-foreground sm:text-lg">
            Frontend: React 18, TypeScript, Redux Toolkit, Tailwind CSS, React Router DOM, i18n for
            internationalization, and comprehensive testing with Jest, React Testing Library, and Cypress. Development
            tools include Storybook, ESLint, Prettier, and Husky for git hooks.
          </p>
          <p className="text-base text-muted-foreground sm:text-lg">
            Backend: Multiple implementations showcasing different approaches - NestJS for a TypeScript-based solution
            with dependency injection and modular architecture, and Gin for a Go-based implementation focusing on
            performance and simplicity. Both follow clean architecture principles and RESTful API design.
          </p>
        </div>
      </section>

      {/* Final CTA Section */}
      <section className="bg-gradient-to-r from-primary via-primary/80 to-primary/60 py-12 text-center sm:py-20 md:py-32">
        <div className="container mx-auto px-4 sm:px-6">
          <h2 className="mb-6 text-3xl font-extrabold text-primary-foreground sm:text-4xl md:text-5xl">
            Explore the Project
          </h2>
          <p className="mx-auto mb-10 max-w-2xl text-base text-primary-foreground/90 sm:text-xl">
            Dive into the codebase to see how modern web development practices are applied in a real-world application,
            featuring dual backend implementations and a robust frontend architecture.
          </p>
          <Button
            asChild
            variant="secondary"
            size="lg"
            className="text-lg font-bold shadow-2xl transition duration-300 ease-in-out hover:scale-110 hover:bg-background/90 focus:outline-none focus:ring-4 focus:ring-background/50"
          >
            <a
              href="https://github.com/singhAmandeep007/ReMinder"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Github
                className="-mt-1 mr-3 inline-block"
                size={24}
              />
              View on GitHub
            </a>
          </Button>
        </div>
      </section>
    </div>
  );
};
