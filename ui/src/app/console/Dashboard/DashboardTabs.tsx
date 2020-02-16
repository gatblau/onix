import React from "react";
import { Tab, Tabs } from "@patternfly/react-core";
import { AddressBookIcon } from "@patternfly/react-icons";

const DashboardTabs: React.FunctionComponent<{}> = () => {

  const [activeTabKey, setActiveTabKey] = React.useState(0);
  // Toggle currently active tab
  const handleTabClick = (event, tabIndex) => {
    setActiveTabKey(tabIndex);
    console.log("====>", tabIndex);
  };

  return (
    <Tabs activeKey={activeTabKey} onSelect={handleTabClick}>
      <Tab eventKey={0} title="HOME">
      </Tab>
      <Tab eventKey={1} title="Tab item 2">
      </Tab>
      <Tab eventKey={2} title="Tab item 3">
      </Tab>
      <Tab
        eventKey={3}
        title={
          <>
            Tab item 4 <AddressBookIcon/>
          </>
        }
      >
      </Tab>
    </Tabs>
  );
};

export default DashboardTabs;
