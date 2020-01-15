import * as React from "react";
import { Route, RouteComponentProps, Switch } from "react-router-dom";
import { accessibleRouteChangeHandler } from "../utils/utils";
import { useDocumentTitle } from "../utils/useDocumentTitle";
import { LastLocationProvider, useLastLocation } from "react-router-last-location";

export interface IAppRoute {
  label?: string;
  /* eslint-disable @typescript-eslint/no-explicit-any */
  component: React.ComponentType<RouteComponentProps<any>> | React.ComponentType<any>;
  exact?: boolean;
  path: string;
  title: string;
  isAsync?: boolean;
}

let routeFocusTimer: number;

// a custom hook for sending focus to the primary content container after a view has loaded so that subsequent press of tab key
// sends focus directly to relevant content
const useA11yRouteChange = (isAsync: boolean) => {
  const lastNavigation = useLastLocation();
  React.useEffect(() => {
    if (!isAsync && lastNavigation !== null) {
      routeFocusTimer = accessibleRouteChangeHandler();
    }
    return () => {
      clearTimeout(routeFocusTimer);
    };
  }, [isAsync, lastNavigation]);
};

const RouteWithTitleUpdates = (
  {
    component: Component,
    isAsync = false,
    title,
    ...rest
  }: IAppRoute) => {
  useA11yRouteChange(isAsync);
  useDocumentTitle(title);

  const routeWithTitle = (routeProps: RouteComponentProps) => (
    <Component {...rest} {...routeProps} />
  );

  return <Route render={routeWithTitle}/>;
};

const ConsoleRoutes = (props) => {
  return (
    <LastLocationProvider>
      <Switch>
        {props.routes.map(({path, exact, component, title, isAsync}, idx) => (
          <RouteWithTitleUpdates
            path={props.baseurl + "/" + path}
            exact={exact}
            component={component}
            key={idx}
            title={title}
            isAsync={isAsync}
          />
        ))}
      </Switch>
    </LastLocationProvider>
  );
};

export default ConsoleRoutes;
