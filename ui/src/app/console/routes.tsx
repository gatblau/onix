import * as React from "react";
import { Dashboard } from "./Dashboard/Dashboard";
import MetaModel from "./MetaModel";
import { IAppRoute } from "./ConsoleLayout/ConsoleRoutes";

const dashboardRoutes: IAppRoute[] = [
  {
    component: Dashboard,
    exact: true,
    label: "Dashboard",
    path: "Dashboard",
    title: "Onix Dashboard"
  }
];
const metaModelRoutes: IAppRoute[] = [
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-1",
    path: "metamodel-1",
    title: "MetaModel"
  },
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-2",
    path: "metamodel-2",
    title: "some other"
  },
  {
    component: MetaModel,
    exact: true,
    label: "MetaModels-3",
    path: "metamodel",
    title: "another"
  }
];

const routes = [...dashboardRoutes, ...metaModelRoutes];

export { routes, dashboardRoutes, metaModelRoutes };
