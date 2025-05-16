import { FC } from "react";

import { Bell } from "lucide-react";

import { useTranslation } from "react-i18next";

import { RouteLink } from "components";

import { HOME_ROUTE_BY_PATH, ROUTE_BY_PATH } from "app/Router";

import { ThemeToggler } from "modules/theme";

import { LangToggler } from "modules/i18n";

export type THeaderProps = Record<string, never>;

export const Header: FC<THeaderProps> = () => {
  const { t } = useTranslation("common");

  return (
    <header
      className="h-[--navbar-height]"
      data-testid="header"
    >
      <div className="flex h-full w-full items-center justify-between border-b-2 border-primary px-4">
        <div>
          <RouteLink to={ROUTE_BY_PATH.home}>
            <Bell className="icon inline-block fill-primary" />
          </RouteLink>
        </div>
        <div>
          <nav className="flex items-center justify-between gap-3">
            {[
              { title: t("navbar.home"), path: ROUTE_BY_PATH.home },
              { title: t("navbar.reminders"), path: HOME_ROUTE_BY_PATH.reminders },
            ].map(({ path, title }) => (
              <RouteLink
                key={path}
                to={path}
              >
                {({ isActive }) => {
                  return <span className={`text-lg ${isActive ? "border-primary text-primary" : ""}`}>{title}</span>;
                }}
              </RouteLink>
            ))}
          </nav>
        </div>
        <div className="inline-flex items-center gap-1">
          <ThemeToggler />
          <LangToggler />
        </div>
      </div>
    </header>
  );
};
