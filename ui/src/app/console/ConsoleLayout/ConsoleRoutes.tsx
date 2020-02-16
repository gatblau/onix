import * as React from "react";
import { Route, RouteComponentProps, Switch } from "react-router-dom";
import { LastLocationProvider } from "react-router-last-location";

export interface IAppRoute {
  label?: string;
  component: React.ComponentType<RouteComponentProps<any>> | React.ComponentType<any>;
  exact?: boolean;
  path: string;
  link: string;
  title: string;
  isAsync?: boolean;
}

const ConsoleRoutes: React.FunctionComponent<any> = (props) => {
  return (
    <LastLocationProvider>
      <Switch>
        {props.routes.map(({path, exact, component}, idx) => (
          <Route
            path={path}
            exact={exact}
            component={component}
            key={idx}
          />
        ))}
      </Switch>
    </LastLocationProvider>
  );
};

export default ConsoleRoutes;
