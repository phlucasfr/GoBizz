import SignUp from "./SignUp";

import { Motion } from "@motionone/solid";
import { useAuth } from "./context/AuthContext";
import { ArrowRight } from "lucide-solid";
import { useNavigate } from "@solidjs/router";
import { validateSession } from "../api/api";
import { getStoredCompanyId } from "../util/Index";
import { createEffect, createSignal } from "solid-js";

const Welcome = () => {
  const id = getStoredCompanyId();
  const navigate = useNavigate();
  const { login } = useAuth();
  const [mounted, setMounted] = createSignal(false);

  createEffect(async () => {
    const sessionData = await validateSession();
    if (sessionData.isValid) {
      login();
      return navigate(`/home/${id}`, { replace: true });
    }
    setMounted(true);
  });

  return (
    <div class="flex flex-col lg:flex-row min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <Motion.div
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: mounted() ? 1 : 0, x: mounted() ? 0 : -20 }}
        transition={{ duration: 0.5 }}
        class="flex flex-col items-center justify-center w-full lg:w-3/5 p-8 lg:p-16"
      >
        <div class="text-center max-w-2xl">
          <h1 class="text-3xl lg:text-5xl font-bold mb-6 text-indigo-900 leading-tight">
            Transforme Sua Gestão Empresarial
          </h1>
          <p class="text-lg lg:text-xl text-indigo-700 mb-8">
            Descubra ferramentas inteligentes e recursos completos para
            impulsionar o sucesso do seu negócio.
          </p>
          <Motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            class="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-full text-lg font-semibold transition-colors duration-300 hover:bg-indigo-700"
          >
            Conheça Nossos Produtos
            <ArrowRight class="ml-2 h-5 w-5" />
          </Motion.button>
        </div>
      </Motion.div>

      <Motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: mounted() ? 1 : 0, x: mounted() ? 0 : 20 }}
        transition={{ duration: 0.5, delay: 0.2 }}
        class="w-screen lg:w-2/5 flex items-center justify-center p-8 lg:p-16"
      >
        <div class="w-screen max-w-md">
          <SignUp />
        </div>
      </Motion.div>
    </div>
  );
};

export default Welcome;
