import * as React from "react";
import { useState } from "react";
import { Avatar, Brand, Nav, NavExpandable, NavItem, NavItemSeparator, NavList, NavVariants, Page, PageHeader, PageSidebar, SkipToContent } from "@patternfly/react-core";
import logo from "assets/logo-text_32.png";
import imgAvatar from 'assets/avatar.svg';
import ConsoleToolbar from "./consoleToolbar";
import * as Routes from "../routes";
import { NavLink, Redirect } from "react-router-dom";

interface IAppLayout {
  children: React.ReactNode;
}

const ConsoleLayout: React.FunctionComponent<IAppLayout> = ({children}) => {
  const logoProps = {
    href: "https://github.com/gatblau/onix",
    target: "_blank"
  };
  const [isNavOpen, setIsNavOpen] = React.useState(true);
  const [isMobileView, setIsMobileView] = React.useState(true);
  const [isNavOpenMobile, setIsNavOpenMobile] = React.useState(false);
  const onNavToggleMobile = () => {
    setIsNavOpenMobile(!isNavOpenMobile);
  };
  const onNavToggle = () => {
    setIsNavOpen(!isNavOpen);
  };
  const onPageResize = (props: { mobileView: boolean; windowSize: number }) => {
    setIsMobileView(props.mobileView);
  };
  const Header = (
    <PageHeader logo={<Brand src={logo} alt="Onix Logo"/>}
                logoProps={logoProps}
                toolbar={ConsoleToolbar}
                avatar={<Avatar src={imgAvatar} alt="TODO User" />}
                isNavOpen={isNavOpen}
                onNavToggle={isMobileView ? onNavToggleMobile : onNavToggle}
                showNavToggle/>
  );
  const [activeItem, setActiveItem] = useState(null);
  const [activeGroup, setActiveGroup] = useState(null);
  const onSelect = result => {
    setActiveGroup(result.groupId);
    setActiveItem(result.itemId);
  };
  const onToggle = result => {
    // eslint-disable-next-line no-console
    console.log(`Group ${result.groupId} expanded? ${result.isExpanded}`);
  };
  const Navigation = (
    <Nav id="onix-console" onSelect={onSelect} onToggle={onToggle} theme="dark">
      <NavList id="onix-console" variant={NavVariants.default}>
        {Routes.dashboardRoutes.map((route, idx) => (
          <NavItem key={`${route.label}-${idx}`} id={`${route.label}-${idx}`}>
            <NavLink exact to={route.path} activeClassName="pf-m-current">{route.title}</NavLink>
          </NavItem>
        ))}
        <NavItemSeparator/>
        <NavExpandable title="Metamodels" groupId="metamodels" isActive={activeGroup === "metamodels"}>
          {Routes.metaModelRoutes.map((route, idx) => (
            <NavItem key={`${route.label}-${idx}`}
                     id={`${route.label}-${idx}`}
                     groupId="metamodels"
                     itemId={`metamodels-${idx}`}
                     isActive={activeItem === `metamodels-${idx}`}>
              <NavLink exact to={route.path}>{route.title}</NavLink>
            </NavItem>
          ))}
        </NavExpandable>
      </NavList>
    </Nav>);
  const Sidebar = (
    <PageSidebar
      theme="dark"
      nav={Navigation}
      isNavOpen={isMobileView ? isNavOpenMobile : isNavOpen}/>
  );
  const PageSkipToContent = (
    <SkipToContent href="#primary-app-container">
      Skip to Content
    </SkipToContent>
  );
  return (
    <Page mainContainerId="primary-app-container"
          header={Header}
          sidebar={Sidebar}
          onPageResize={onPageResize}
          skipToContent={PageSkipToContent}>
      {children}
    </Page>
  );
};

export { ConsoleLayout };
