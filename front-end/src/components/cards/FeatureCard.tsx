import { Motion } from "@motionone/solid";
import { BarChart3, FileText, Link, Package, ShoppingCart, Users } from "lucide-solid";

export type FeatureCardProps = {
  href: string;
  icon: typeof BarChart3;
  color: string;
  title: string;
};

export const features: FeatureCardProps[] = [
  {
    title: "Cadastro de Funcionários",
    icon: Users,
    color: "text-blue-500",
    href: "/employees",
  },
  {
    title: "E-commerce",
    icon: ShoppingCart,
    color: "text-green-500",
    href: "/ecommerce",
  },
  {
    title: "Gerador de Relatórios",
    icon: FileText,
    color: "text-yellow-500",
    href: "/reports",
  },
  {
    title: "Controle de Estoque",
    icon: Package,
    color: "text-purple-500",
    href: "/inventory",
  },
  {
    title: "Encurtador de Links",
    icon: Link,
    color: "text-red-500",
    href: "/links",
  },
];

const FeatureCard = (props: FeatureCardProps) => (
  <Motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
    <div class="bg-white p-6 rounded-lg shadow-md cursor-pointer hover:shadow-lg transition-shadow duration-300">
      <a href={props.href}>
        <props.icon class={`h-8 w-8 ${props.color}`} />
        <h3 class="mt-4 text-lg font-semibold">{props.title}</h3>
        <p class="mt-2 text-sm text-gray-600">Clique para acessar</p>
      </a>
    </div>
  </Motion.div>
);

export default FeatureCard;
