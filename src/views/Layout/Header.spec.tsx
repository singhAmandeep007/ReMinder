import { screen } from "@testing-library/react";

import { render } from "tests/utils";

import { Header } from "./Header";

describe("Header", () => {
  const setup = () => {
    return {
      result: render(<Header />),
    };
  };

  it("should render links", () => {
    setup();

    expect(screen.getByTestId("header")).toBeInTheDocument();

    ["ReMinder", "Explore Demo", "Github"].forEach((link) => {
      expect(screen.getByRole("link", { name: `${link}` })).toBeInTheDocument();
    });
  });
});
