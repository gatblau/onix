import * as React from "react";
import DataGraph from "./DataGraph";
import DataTable from "./DataTable";
import { PageSection, Title } from "@patternfly/react-core";

const MetaModel: React.FunctionComponent<{}> = (props) => {
  return (
    <PageSection>
      <Title size="lg">MetaModel</Title>
      <DataGraph/>
      <DataTable/>
    </PageSection>);
};

export default MetaModel;
