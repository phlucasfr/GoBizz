import Sidebar from "./Sidebar";
import LogoutModal from "./modals/LogoutModal";
import QuickStatsCard from "./cards/QuickStatsCard";

import { Motion } from "@motionone/solid";
import { useAuth } from "./context/AuthContext";
import { useParams } from "@solidjs/router";
import { useNavigate } from "@solidjs/router";
import { deleteSession, validateSession } from "../api/api";
import FeatureCard, { features } from "./cards/FeatureCard";
import { createSignal, createEffect, For, Show } from "solid-js";

const Home = () => {
  const navigate = useNavigate();
  const { login, logout } = useAuth();

  const params = useParams();
  const [userName, setUserName] = createSignal("");
  const [isLogoutModalOpen, setIsLogoutModalOpen] = createSignal(false);

  const handleLogout = async () => {
    const response = await deleteSession();
    if (!response.success) return;
    
    logout()
    setIsLogoutModalOpen(false);
    navigate("/", { replace: true });
  };

  createEffect(async () => {
    const sessionData = await validateSession();
    if (!sessionData.isValid) {
      logout()
      return navigate(`/login`, { replace: true });
    }
    
    login();
  });

  createEffect(() => {
    setUserName("Usuário");
    //TODO: Set Company Name and other infos.
    console.log("User ID:", params.id);
  });

  return (
    <div class="flex h-screen bg-gray-100">
      <Sidebar onLogoutClick={() => setIsLogoutModalOpen(true)} />
      <main class="flex-1 overflow-y-auto p-8">
        <Motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          class="max-w-7xl mx-auto"
        >
          <h1 class="text-3xl font-bold text-gray-900 mb-8">
            Bem-vindo, {userName()}!
          </h1>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            <QuickStatsCard
              title="Vendas Hoje"
              value="R$ 5.230"
              change="+12%"
            />
            <QuickStatsCard title="Novos Clientes" value="24" change="+8%" />
            <QuickStatsCard
              title="Itens em Estoque"
              value="1.423"
              change="-3%"
            />
          </div>

          <h2 class="text-2xl font-semibold text-gray-800 mb-4">
            Funcionalidades
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <For each={features}>
              {(feature) => (
                <FeatureCard
                  title={feature.title}
                  icon={feature.icon}
                  color={feature.color}
                  href={feature.href}
                />
              )}
            </For>
          </div>
        </Motion.div>
      </main>

      <Show when={isLogoutModalOpen()}>
        <LogoutModal
          isOpen={isLogoutModalOpen()}
          onClose={() => setIsLogoutModalOpen(false)}
          onLogout={handleLogout}
        />
      </Show>
    </div>
  );
};

export default Home;
