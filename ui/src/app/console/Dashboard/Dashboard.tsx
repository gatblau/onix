import * as React from "react";
import { PageSection, Title } from "@patternfly/react-core";
import DashboardTabs from "./DashboardTabs";

const Dashboard: React.FunctionComponent<{}> = () => (
  <PageSection>
    <Title size="lg">Dashboard</Title>
    <DashboardTabs/>
  </PageSection>
);

export { Dashboard };
