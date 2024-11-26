import SmsVerificationModal from "./SmsVerificationModal";

import { useNavigate } from "@solidjs/router";
import { createSignal, createResource, Show } from "solid-js";
import { maskCpfCnpj, maskPhone, validateEmail } from "../util/Index";
import { apiConfig } from "../config/apiConfig";

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

    if (!name || !email || !phone || !cpfCnpj || password.length < 8) {
      throw new Error("Preencha todos os campos obrigatórios corretamente.");
    }

    if (!validateEmail(email)) {
      throw new Error("E-mail inválido.");
    }

    const payload = {
      name,
      email,
      phone: cleanPhone,
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

      <div class="min-h-screen flex items-center justify-center p-4">
        <div class="w-[50vh] max-w-md p-8 rounded-lg shadow-lg space-y-6">
          <h2 class="text-3xl font-semibold text-blue-600 text-center">
            Cadastre Sua Empresa
          </h2>

          <Show when={error()}>
            <div class="mb-4 p-4 bg-red-100 text-red-800 border-l-4 border-red-500 rounded-lg">
              {error()}
            </div>
          </Show>

          <form onSubmit={handleSubmit} class="space-y-6">
            <div>
              <label for="name" class="block text-lg font-medium text-gray-700">
                Nome da Empresa
              </label>
              <input
                id="name"
                name="name"
                type="text"
                value={formData().name}
                onInput={handleChange}
                required
                class="text-black mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label
                for="email"
                class="block text-lg font-medium text-gray-700"
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
                class="text-black mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label
                for="phone"
                class="block text-lg font-medium text-gray-700"
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
                class="text-black mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label
                for="cpfCnpj"
                class="block text-lg font-medium text-gray-700"
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
                class="text-black mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label
                for="password"
                class="block text-lg font-medium text-gray-700"
              >
                Senha
              </label>
              <input
                id="password"
                name="password"
                type="password"
                value={formData().password}
                onInput={handleChange}
                required
                class="text-black mt-2 block w-full px-5 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <button
              type="submit"
              disabled={data.loading}
              class={`w-full py-3 px-5 rounded-lg shadow-md text-white ${
                data.loading ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-500"
              } focus:outline-none focus:ring-2 focus:ring-blue-600`}
            >
              {data.loading ? "Carregando..." : "Cadastrar"}
            </button>
          </form>

          <div class="mt-4 text-center text-black">
            <p class="text-sm">
              Já tem uma conta?{" "}
              <a
                onClick={() => navigate("/login")}
                class="text-blue-600 hover:text-blue-500 cursor-pointer"
              >
                Faça login
              </a>
            </p>
          </div>
        </div>
      </div>
    </>
  );
};

export default SignUp;
