import { FC } from "react";
import { Outlet } from "react-router";

import { Header } from "./Header";

import { Footer } from "./Footer";

export type TBaseLayoutProps = Record<string, never>;

export const BaseLayout: FC<TBaseLayoutProps> = () => {
  return (
    <div>
      <Header />

      <main
        className="min-h-[--content-height] overflow-x-auto"
        data-testid="content"
      >
        <Outlet />
      </main>

      <Footer />
    </div>
  );
};
