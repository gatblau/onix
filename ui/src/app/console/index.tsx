import * as React from "react";
import "@patternfly/react-core/dist/styles/base.css";
import { BrowserRouter as Router, Redirect, Switch, withRouter } from "react-router-dom";

import ConsoleLayout from "./ConsoleLayout/ConsoleLayout";
import ConsoleRoutes from "./ConsoleLayout/ConsoleRoutes";
import { routes } from "./routes";

const Console: React.ComponentClass<{}> = withRouter((props) => {
  return (
    <Router>
      <ConsoleLayout>
        <Switch>
          <ConsoleRoutes routes={routes} baseurl={props.match.url}/>
        </Switch>
      </ConsoleLayout>
      <Redirect to={"/console/dashboard"}/>
    </Router>);
});

export default Console;
