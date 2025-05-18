import { screen } from "@testing-library/react";

import { render } from "tests/utils";

import { Home } from "./Home";

describe("Home", () => {
  it("should render correctly", () => {
    render(<Home />);

    expect(screen.getByText("Full-Stack Reminder Application")).toBeInTheDocument();
    expect(screen.getByText("Explore Demo")).toBeInTheDocument();
    expect(screen.getByText("Check Build Process")).toBeInTheDocument();
  });
});
