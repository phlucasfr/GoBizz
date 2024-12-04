import PasswordRecoveryModal from "./modals/PasswordRecoveryModal";

import { Motion } from "@motionone/solid";
import { useAuth } from "./context/AuthContext";
import { useNavigate } from "@solidjs/router";
import { Eye, EyeOff } from "lucide-solid";
import { validateEmail } from "../util/Index";
import { loginCompany, validateSession } from "../api/api";
import { createSignal, createEffect, Show } from "solid-js";

const SignIn = () => {
  const id = localStorage.getItem("id");
  const navigate = useNavigate();
  const { login } = useAuth();

  const [error, setError] = createSignal("");
  const [mounted, setMounted] = createSignal(false);
  const [formData, setFormData] = createSignal({
    email: "",
    password: "",
  });
  const [isRecoveryModalOpen, setRecoveryModalOpen] = createSignal(false);
  const [showPassword, setShowPassword] = createSignal(false);
  const [isSubmitting, setIsSubmitting] = createSignal(false);

  createEffect(async () => {
    const sessionData = await validateSession();

    if (sessionData.isValid) {
      login();
      return navigate(`/home/${id}`, { replace: true });
    }
    setMounted(true);
  });

  createEffect(() => {
    if (error()) {
      const errorElement = document.getElementById("error-message");
      errorElement?.focus();
    }
  });

  const handleChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    const value = target.type === "checkbox" ? target.checked : target.value;

    setFormData({ ...formData(), [target.name]: value });
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    const { email, password } = formData();

    setError("");
    setIsSubmitting(true);

    if (!email || password.length < 8) {
      setError("Por favor, preencha todos os campos corretamente.");
      setIsSubmitting(false);
      return;
    }

    if (!validateEmail(email)) {
      setError("Por favor, insira um e-mail válido.");
      setIsSubmitting(false);
      return;
    }

    const response = await loginCompany(email, password);

    if (!response.success) {
      setError(response.message);
      setIsSubmitting(false);
      return;
    }

    login();
    navigate(`/home/${id}`, { replace: true });
  };

  return (
    <div class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-blue-50 to-indigo-100">
      <Motion.div
        initial={{ opacity: 0, x: 20, scale: 0.9 }}
        animate={{
          opacity: mounted() ? 1 : 0,
          x: mounted() ? 0 : 20,
          scale: mounted() ? 1 : 0.9,
        }}
        transition={{ duration: 0.5 }}
        class="w-full max-w-md p-8 bg-white rounded-xl shadow-lg space-y-6"
      >
        <Motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <img
            class="h-12 w-auto mx-auto mb-6"
            src="/assets/logo.png"
            alt="GoBizz Logo"
          />
        </Motion.div>

        <Motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <h1 class="text-3xl font-bold text-center text-indigo-900 mb-6">
            Bem-vindo de volta
          </h1>
        </Motion.div>

        <Show when={error()}>
          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.3 }}
            id="error-message"
            role="alert"
            tabIndex={-1}
            class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg"
          >
            {error()}
          </Motion.div>
        </Show>

        <form onSubmit={handleSubmit} class="space-y-6">
          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.4 }}
          >
            <label
              for="email"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              E-mail
            </label>
            <input
              id="email"
              name="email"
              type="email"
              value={formData().email}
              onInput={handleChange}
              required
              class="block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              placeholder="seu@email.com"
            />
          </Motion.div>

          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.5 }}
            class="relative"
          >
            <label
              for="password"
              class="block text-sm font-medium text-gray-700 mb-1"
            >
              Senha
            </label>
            <input
              id="password"
              name="password"
              type={showPassword() ? "text" : "password"}
              value={formData().password}
              onInput={handleChange}
              required
              class="block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              placeholder="Sua senha"
            />
            <Motion.div
              class="absolute right-3 top-8 text-gray-500 focus:outline-none"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
            >
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword())}
                aria-label={showPassword() ? "Esconder senha" : "Mostrar senha"}
              >
                {showPassword() ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </Motion.div>
          </Motion.div>

          <Motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.6 }}
            class="flex items-center justify-between"
          >
            <Motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <a
                href="#"
                class="font-medium text-indigo-600 hover:text-indigo-500 transition-colors duration-200"
                onclick={setRecoveryModalOpen}
              >
                Esqueceu a senha?
              </a>
            </Motion.div>
          </Motion.div>

          <Motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.7 }}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <button
              type="submit"
              disabled={isSubmitting()}
              class="w-full bg-gradient-to-r from-indigo-600 to-blue-500 text-white py-3 px-5 rounded-lg shadow-md hover:from-indigo-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            >
              {isSubmitting() ? "Entrando..." : "Entrar"}
            </button>
          </Motion.div>
        </form>

        <Motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.8 }}
          class="mt-4 text-center text-gray-600"
        >
          <p class="text-sm">
            Não tem uma conta?{" "}
            <Motion.div
              class="inline-block"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <a
                onClick={() => navigate("/register")}
                class="text-indigo-600 hover:text-indigo-500 cursor-pointer font-medium transition-colors duration-200"
              >
                Registre-se
              </a>
            </Motion.div>
          </p>
        </Motion.div>
      </Motion.div>

      <PasswordRecoveryModal
        isOpen={isRecoveryModalOpen()}
        onClose={() => setRecoveryModalOpen(false)}
      ></PasswordRecoveryModal>
    </div>
  );
};

export default SignIn;
