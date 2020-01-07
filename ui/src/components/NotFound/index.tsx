import * as React from "react";
import { NavLink } from "react-router-dom";
import { Alert, PageSection } from "@patternfly/react-core";

const NotFound: React.FunctionComponent = () => (
  <PageSection>
    <Alert variant="danger" title="404! I haven't build this yet :)"/><br/>
    <NavLink to="/" className="pf-c-nav__link">Take me home</NavLink>
  </PageSection>
);

export { NotFound };
