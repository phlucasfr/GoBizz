import SmsVerificationModal from "./SmsVerificationModal";
import { Motion } from "@motionone/solid";
import { apiConfig } from "../config/apiConfig";
import { useNavigate } from "@solidjs/router";
import { Eye, EyeOff } from 'lucide-solid';
import { createSignal, createResource, Show } from "solid-js";
import { maskCpfCnpj, maskPhone, validateEmail } from "../util/Index";

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
  const [showPassword, setShowPassword] = createSignal(false);
  const [isModalVisible, setIsModalVisible] = createSignal(false);

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
    const lowerCaseEmail = email.toLowerCase()

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

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || "Erro ao cadastrar a empresa.");
    }

    const data = await response.json();
    localStorage.setItem("id", data.id);

    setIsModalVisible(true);
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
      <SmsVerificationModal
        isVisible={isModalVisible()}
        onClose={() => {
          setIsModalVisible(false);
        }}
      />

      <div class="min-h-fit flex items-center justify-center p-4 ">
        <Motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          class="w-full max-w-md p-8 bg-white rounded-xl shadow-2xl space-y-6"
        >
          <h2 class="text-3xl font-bold text-indigo-700 text-center">
            Cadastre Sua Empresa
          </h2>

          <Show when={error()}>
            <Motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.3 }}
              class="p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg"
            >
              {error()}
            </Motion.div>
          </Show>

          <form onSubmit={handleSubmit} class="space-y-6">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700">
                Nome da Empresa
              </label>
              <input
                id="name"
                name="name"
                type="text"
                value={formData().name}
                onInput={handleChange}
                required
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
                       focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
                       disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
                       invalid:border-pink-500 invalid:text-pink-600
                       focus:invalid:border-pink-500 focus:invalid:ring-pink-500"
              />
            </div>

            <div>
              <label for="email" class="block text-sm font-medium text-gray-700">
                E-mail
              </label>
              <input
                id="email"
                name="email"
                type="email"
                value={formData().email}
                onInput={handleChange}
                required
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
                       focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
                       disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
                       invalid:border-pink-500 invalid:text-pink-600
                       focus:invalid:border-pink-500 focus:invalid:ring-pink-500"
              />
            </div>

            <div>
              <label for="phone" class="block text-sm font-medium text-gray-700">
                Telefone
              </label>
              <input
                id="phone"
                name="phone"
                type="text"
                value={formData().phone}
                onInput={handleChange}
                required
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
                       focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
                       disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
                       invalid:border-pink-500 invalid:text-pink-600
                       focus:invalid:border-pink-500 focus:invalid:ring-pink-500"
              />
            </div>

            <div>
              <label for="cpfCnpj" class="block text-sm font-medium text-gray-700">
                CPF/CNPJ
              </label>
              <input
                id="cpfCnpj"
                name="cpfCnpj"
                type="text"
                value={formData().cpfCnpj}
                onInput={handleChange}
                required
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md text-sm shadow-sm placeholder-gray-400
                       focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500
                       disabled:bg-gray-50 disabled:text-gray-500 disabled:border-gray-200 disabled:shadow-none
                       invalid:border-pink-500 invalid:text-pink-600
                       focus:invalid:border-pink-500 focus:invalid:ring-pink-500"
              />
            </div>

            <div>
              <label for="password" class="block text-sm font-medium text-gray-700">
                Senha
              </label>
              <div class="mt-1 relative rounded-md shadow-sm">
                <input
                  id="password"
                  name="password"
                  type={showPassword() ? "text" : "password"}
                  value={formData().password}
                  onInput={handleChange}
                  required
                  class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400
                         focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                />
                <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword())}
                    class="text-gray-400 hover:text-gray-500 focus:outline-none focus:text-gray-500"
                  >
                    {showPassword() ? <EyeOff class="h-5 w-5" /> : <Eye class="h-5 w-5" />}
                  </button>
                </div>
              </div>
            </div>

            <Motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <button
                type="submit"
                disabled={data.loading}
                class={`w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white 
                        ${data.loading ? "bg-indigo-400 cursor-not-allowed" : "bg-indigo-600 hover:bg-indigo-700"}
                        focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500`}
              >
                {data.loading ? "Carregando..." : "Cadastrar"}
              </button>
            </Motion.div>
          </form>

          <div class="mt-4 text-center text-sm text-gray-600">
            <p>
              Já tem uma conta?{" "}
              <a
                onClick={() => navigate("/login")}
                class="font-medium text-indigo-600 hover:text-indigo-500 cursor-pointer"
              >
                Faça login
              </a>
            </p>
          </div>
        </Motion.div>
      </div>
    </>
  );
};

export default SignUp;