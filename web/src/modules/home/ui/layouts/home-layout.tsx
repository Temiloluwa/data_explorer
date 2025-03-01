import { SidebarProvider } from "@/components/ui/sidebar";
import { HomeNavBar } from "@/modules/home/ui/components/home-navbar";
import { HomeSideBar } from "@/modules/home/ui/components/home-sidebar";
import { HomeFooter } from "@/modules/home/ui/components/home-footer";

interface HomeLayoutProps {
  children: React.ReactNode;
}

export const HomeLayout = ({ children }: HomeLayoutProps) => {
  return (
    <SidebarProvider>
      <div className="w-full flex">
        <HomeSideBar />
        <main className="flex flex-col justify-between min-h-screen w-full">
          <HomeNavBar />
          {children}
          <HomeFooter />
        </main>
      </div>
    </SidebarProvider>
  );
};
