import * as React from "react";
import "@patternfly/react-core/dist/styles/base.css";
import { BrowserRouter as Router } from "react-router-dom";

import { ConsoleLayout } from "./ConsoleLayout/ConsoleLayout";
import { ConsoleRoutes } from "./routes";

const Console: React.FunctionComponent<{}> = () => {
  return (
    <Router>
      <ConsoleLayout>
        <ConsoleRoutes/>
      </ConsoleLayout>
    </Router>);
};

export default Console;
