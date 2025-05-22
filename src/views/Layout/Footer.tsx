import { FC } from "react";

import { useTranslation } from "react-i18next";

import { Typography } from "components";

export type TFooterProps = Record<string, never>;

export const Footer: FC<TFooterProps> = () => {
  const { t } = useTranslation("common");
  return (
    <footer
      className="h-min-[--footer-min-height]"
      data-testid="footer"
    >
      <div className="flex h-full w-full items-center justify-center gap-2 border-t-2 border-primary p-2">
        <Typography variant={"p"}>
          {t("app.copyright", {
            year: new Date().getFullYear(),
          })}
          <a
            target="_blank"
            rel="noreferrer"
            href="https://singhamandeep007.github.io"
            className="ml-1 border-primary text-primary hover:border-b-2"
          >
            {t("app.author")}
          </a>
        </Typography>

        <span>&bull;</span>

        <nav className="flex flex-wrap items-center justify-center gap-4">
          <a
            href="https://github.com/singhAmandeep007/ReMinder"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-primary"
          >
            GitHub
          </a>
          <span>&bull;</span>
          <a
            href="https://www.linkedin.com/in/singhamandeep007/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-primary"
          >
            LinkedIn
          </a>
        </nav>
      </div>
    </footer>
  );
};
