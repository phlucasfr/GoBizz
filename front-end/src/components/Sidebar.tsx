import { A } from "@solidjs/router";
import { For } from "solid-js";
import {
  Home,
  Users,
  LogOut,
  Package,
  Settings,
  FileText,
  ShoppingCart,
  Link as LinkIcon,
} from "lucide-solid";

type MenuItem = {
  href: string;
  icon: typeof Home;
  label: string;
  action?: () => void;
};

const Sidebar = (props: { onLogoutClick: () => void }) => {
  const menuItems: MenuItem[] = [
    { icon: Home, label: "Dashboard", href: "/" },
    { icon: Users, label: "Funcionários", href: "/employees" },
    { icon: ShoppingCart, label: "E-commerce", href: "/ecommerce" },
    { icon: FileText, label: "Relatórios", href: "/reports" },
    { icon: Package, label: "Estoque", href: "/inventory" },
    { icon: LinkIcon, label: "Links", href: "/links" },
    { icon: Settings, label: "Configurações", href: "/settings" },
    {
      icon: LogOut,
      label: "Sair",
      href: "#",
      action: props.onLogoutClick,
    },
  ];

  return (
    <aside class="bg-white w-64 min-h-screen p-4 border-r border-gray-200">
      <div class="flex items-center justify-center mb-8">
        <img src="/assets/logo.png" alt="Logo" class="h-8 w-auto" />
        <span class="ml-2 text-xl font-bold text-gray-800">GoBizz</span>
      </div>
      <nav>
        <ul class="space-y-2">
          <For each={menuItems}>
            {(item) => (
              <li>
                <A
                  href={item.href}
                  class="flex items-center space-x-3 text-gray-700 p-2 rounded-lg font-medium hover:bg-gray-200 focus:shadow-outline"
                  onClick={(e) => {
                    e.preventDefault();
                    item.action?.();
                  }}
                >
                  <item.icon class="h-5 w-5" />
                  <span>{item.label}</span>
                </A>
              </li>
            )}
          </For>
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;
