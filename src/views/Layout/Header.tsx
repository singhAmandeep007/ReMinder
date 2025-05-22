import { FC } from "react";

import { Code2, Github } from "lucide-react";

import { useTranslation } from "react-i18next";

import { useLocation } from "react-router";

import { cn } from "utils";

import { RouteLink, Button } from "components";

import { HOME_ROUTE_BY_PATH, ROUTE_BY_PATH } from "app/Router";

import { ThemeToggler } from "modules/theme";

import { LangToggler } from "modules/i18n";

export type THeaderProps = Record<string, never>;

export const Header: FC<THeaderProps> = () => {
  const { t } = useTranslation("common");

  const location = useLocation();

  const isHomePath = location.pathname === ROUTE_BY_PATH.home;

  return (
    <header
      className={`${cn(isHomePath ? "sticky" : "")} top-0 z-[--navbar-z-index] border-b-2 border-primary bg-background/80 p-4 shadow-lg backdrop-blur-sm`}
      data-testid="header"
    >
      <div className="mx-auto flex flex-row items-center justify-between gap-4">
        <RouteLink
          to={ROUTE_BY_PATH.home}
          className="flex items-center text-xl font-bold"
        >
          <Code2 className="mr-2 h-8 w-8 text-primary" />
          {t("app.appName")}
        </RouteLink>

        <div className="flex flex-wrap items-center justify-center gap-4">
          {isHomePath && (
            <Button
              asChild
              className="hidden md:inline-flex"
            >
              <RouteLink to={HOME_ROUTE_BY_PATH.reminders}> {t("navbar.exploreDemo")}</RouteLink>
            </Button>
          )}

          <Button
            asChild
            variant="secondary"
            className="hidden md:inline-flex"
          >
            <a
              href="https://github.com/singhAmandeep007/ReMinder"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Github className="mr-2 inline-block h-5 w-5" />
              {t("navbar.github")}
            </a>
          </Button>
          <div className="flex items-center gap-2">
            <ThemeToggler />
            <LangToggler />
          </div>
        </div>
      </div>
    </header>
  );
};
