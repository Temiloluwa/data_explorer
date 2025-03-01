import { Sidebar, SidebarContent } from "@/components/ui/sidebar";
import { SidebarAppsSection } from "@/modules/home/ui/components/home-sidebar/apps-section";
import { SidebarFooterSection } from "@/modules/home/ui/components/home-sidebar/footer-section";
import { SidebarHeaderSection } from "@/modules/home/ui/components/home-sidebar/header-section";

export const HomeSideBar = () => {
  return (
    <Sidebar collapsible="icon">
      <SidebarHeaderSection />
      <SidebarContent>
        <SidebarAppsSection />
      </SidebarContent>
      <SidebarFooterSection />
    </Sidebar>
  );
};
