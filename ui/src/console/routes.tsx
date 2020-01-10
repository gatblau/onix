import * as React from "react";
import { Route, RouteComponentProps, Switch } from "react-router-dom";
import { accessibleRouteChangeHandler } from "utils/utils";
import { Dashboard } from "./Dashboard/Dashboard";

import { useDocumentTitle } from "utils/useDocumentTitle";
import { LastLocationProvider, useLastLocation } from "react-router-last-location";
import MetaModel from "./MetaModel";

let routeFocusTimer: number;

export interface IAppRoute {
  label?: string;
  /* eslint-disable @typescript-eslint/no-explicit-any */
  component: React.ComponentType<RouteComponentProps<any>> | React.ComponentType<any>;
  exact?: boolean;
  path: string;
  title: string;
  isAsync?: boolean;
}

const dashboardRoutes: IAppRoute[] = [
  {
    component: Dashboard,
    exact: true,
    label: "Dashboard",
    path: "/console/Dashboard",
    title: "Onix Dashboard"
  }
];
const metaModelRoutes: IAppRoute[] = [
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-1",
    path: "/console/metamodel-1",
    title: "MetaModel"
  },
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-2",
    path: "/console/metamodel-2",
    title: "some other"
  },
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-3",
    path: "/console/metamodel",
    title: "another"
  }
];

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

const routes = [...dashboardRoutes, ...metaModelRoutes];

const ConsoleRoutes = () => (
  <LastLocationProvider>
    <Switch>
      {routes.map(({path, exact, component, title, isAsync}, idx) => (
        <RouteWithTitleUpdates
          path={path}
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

// TODO make routes dynamic
export { ConsoleRoutes, dashboardRoutes, metaModelRoutes };
