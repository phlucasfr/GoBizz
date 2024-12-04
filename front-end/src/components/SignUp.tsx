import EmailVerificationModal from "./modals/EmailVerificationModal";

import { Motion } from "@motionone/solid";
import { apiConfig } from "../config/apiConfig";
import { useNavigate } from "@solidjs/router";
import { Eye, EyeOff } from "lucide-solid";
import { maskCpfCnpj, maskPhone, validateEmail } from "../util/Index";
import { createSignal, createEffect, createResource, Show } from "solid-js";

const SignUp = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = createSignal({
    name: "",
    email: "",
    phone: "",
    cpfCnpj: "",
    password: "",
  });
  const [error, setError] = createSignal("");
  const [mounted, setMounted] = createSignal(false);
  const [showPassword, setShowPassword] = createSignal(false);
  const [isEmailModalVisible, setEmailModalVisible] = createSignal(false);

  createEffect(() => {
    setMounted(true);
  });

  const handleChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    let value = target.value;

    if (target.name === "phone") {
      value = maskPhone(value);
    } else if (target.name === "cpfCnpj") {
      value = maskCpfCnpj(value);
    }

    setFormData({ ...formData(), [target.name]: value });
  };

  const submitForm = async () => {
    const { name, email, phone, cpfCnpj, password } = formData();

    const cleanPhone = phone.replace(/\D/g, "");
    const cleanCpfCnpj = cpfCnpj.replace(/\D/g, "");
    const lowerCaseEmail = email.toLowerCase();

    if (!name || !email || !phone || !cpfCnpj || password.length < 8) {
      throw new Error("Preencha todos os campos obrigatórios corretamente.");
    }

    if (!validateEmail(email)) {
      throw new Error("E-mail inválido.");
    }

    const payload = {
      name,
      phone: cleanPhone,
      email: lowerCaseEmail,
      cpf_cnpj: cleanCpfCnpj,
      password,
    };

    const response = await fetch(
      `${apiConfig.baseUrl}${apiConfig.endpoints.company}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      }
    );

    const data = await response.json();
    if (!response.ok) return setError(data.error);

    localStorage.setItem("id", data.id);
    setEmailModalVisible(true);
  };

  const [data, { refetch }] = createResource(submitForm);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");

    try {
      await refetch();
    } catch (err: any) {
      setError(err.message || "Erro inesperado. Tente novamente.");
    }
  };

  return (
    <>
      <EmailVerificationModal
        isVisible={isEmailModalVisible()}
        onClose={() => {
          setEmailModalVisible(false);
        }}
      />
      <Motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: mounted() ? 1 : 0 }}
        transition={{ duration: 0.5 }}
        class="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4 sm:p-6 lg:p-8"
      >
        <Motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{
            opacity: mounted() ? 1 : 0,
            scale: mounted() ? 1 : 0.9,
          }}
          transition={{ duration: 0.5 }}
          class="w-full max-w-md bg-white rounded-xl shadow-lg p-6 sm:p-8 space-y-6"
        >
          <Motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <h2 class="text-2xl sm:text-3xl font-bold text-center text-indigo-900 mb-6">
              Cadastre Sua Empresa
            </h2>
          </Motion.div>

          <Show when={error()}>
            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.3 }}
              id="error-message"
              role="alert"
              tabIndex={-1}
              class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg text-sm"
            >
              {error()}
            </Motion.div>
          </Show>

          <form onSubmit={handleSubmit} class="space-y-4 sm:space-y-6">
            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.3 }}
            >
              <label
                for="name"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                Nome da Empresa
              </label>
              <input
                id="name"
                name="name"
                type="text"
                value={formData().name}
                onInput={handleChange}
                required
                class="block w-full px-3 py-2 sm:px-4 sm:py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm sm:text-base"
                placeholder="Sua empresa"
              />
            </Motion.div>

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
                class="block w-full px-3 py-2 sm:px-4 sm:py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm sm:text-base"
                placeholder="seu@email.com"
              />
            </Motion.div>

            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.5 }}
            >
              <label
                for="phone"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                Telefone
              </label>
              <input
                id="phone"
                name="phone"
                type="text"
                value={formData().phone}
                onInput={handleChange}
                required
                class="block w-full px-3 py-2 sm:px-4 sm:py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm sm:text-base"
                placeholder="(00) 00000-0000"
              />
            </Motion.div>

            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.6 }}
            >
              <label
                for="cpfCnpj"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                CPF/CNPJ
              </label>
              <input
                id="cpfCnpj"
                name="cpfCnpj"
                type="text"
                value={formData().cpfCnpj}
                onInput={handleChange}
                required
                class="block w-full px-3 py-2 sm:px-4 sm:py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm sm:text-base"
                placeholder="000.000.000-00 ou 00.000.000/0000-00"
              />
            </Motion.div>

            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.7 }}
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
                class="block w-full px-3 py-2 sm:px-4 sm:py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm sm:text-base"
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
                  aria-label={
                    showPassword() ? "Esconder senha" : "Mostrar senha"
                  }
                >
                  {showPassword() ? <EyeOff size={20} /> : <Eye size={20} />}
                </button>
              </Motion.div>
            </Motion.div>

            <Motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.8 }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <button
                type="submit"
                disabled={data.loading}
                class="w-full bg-gradient-to-r from-indigo-600 to-blue-500 text-white py-2 sm:py-3 px-4 sm:px-5 rounded-lg shadow-md hover:from-indigo-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 text-sm sm:text-base font-medium"
              >
                {data.loading ? "Cadastrando..." : "Cadastrar"}
              </button>
            </Motion.div>
          </form>

          <Motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.9 }}
            class="mt-4 text-center text-gray-600"
          >
            <p class="text-sm">
              Já tem uma conta?{" "}
              <Motion.div
                class="inline-block"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <a
                  onClick={() => navigate("/login")}
                  class="text-indigo-600 hover:text-indigo-500 cursor-pointer font-medium transition-colors duration-200"
                >
                  Faça login
                </a>
              </Motion.div>
            </p>
          </Motion.div>
        </Motion.div>
      </Motion.div>
    </>
  );
};

export default SignUp;
