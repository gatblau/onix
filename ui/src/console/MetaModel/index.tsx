import * as React from "react";
import DataGraph from "./DataGraph";
import DataTable from "./DataTable";
import { PageSection, Title } from "@patternfly/react-core";

const MetaModel: React.FunctionComponent<{}> = () => (
  <PageSection>
    <Title size="lg">MetaModel</Title>
    <DataGraph/>
    <DataTable name="Steve"/>
  </PageSection>
);

export default MetaModel;
