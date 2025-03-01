import { SidebarTrigger } from "@/components/ui/sidebar";

export const HomeNavBar = () => {
  return (
    <nav className="flex relative w-full">
      <div className="absolute top-0 left-0 right-0 flex items-center justify-between w-full h-12 z-50 px-3 bg-slate-300">
        <SidebarTrigger />
        <div className="">
          <p>User Logo</p>
        </div>
      </div>
    </nav>
  );
};
