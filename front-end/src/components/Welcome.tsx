import { Motion } from "solid-motionone";
import { useNavigate } from "@solidjs/router";
import { createEffect, createSignal, For } from "solid-js";
import { FileText, Link, Zap, Shield, ArrowRight } from "lucide-solid";

const products = [
  {
    title: "Gerador de Relatórios",
    icon: FileText,
    color: "text-blue-600",
    bgColor: "bg-blue-100",
    description:
      "Crie relatórios detalhados e insights valiosos para seu negócio.",
    details:
      "Transforme dados brutos em insights acionáveis com nosso gerador de relatórios intuitivo. Crie visualizações impressionantes e relatórios detalhados que ajudam a tomar decisões informadas e impulsionar o crescimento do seu negócio.",
    performance:
      "Processamento rápido de grandes conjuntos de dados para relatórios em tempo real.",
    security:
      "Controle de acesso granular e auditoria de logs para proteger informações sensíveis.",
  },
  {
    title: "Encurtador de Links",
    icon: Link,
    color: "text-green-600",
    bgColor: "bg-green-100",
    description:
      "Simplifique e rastreie seus links com nossa ferramenta intuitiva.",
    details:
      "Crie links curtos e memoráveis para suas campanhas de marketing, produtos ou conteúdo. Nossa ferramenta de encurtamento de links oferece análises detalhadas de cliques, geolocalização e dispositivos, permitindo que você otimize suas estratégias de compartilhamento.",
    performance:
      "Redirecionamento ultrarrápido para uma experiência de usuário sem interrupções.",
    security: "Proteção contra spam e malware para links seguros e confiáveis.",
  },
];

const Welcome = () => {
  const navigate = useNavigate();
  const [mounted, setMounted] = createSignal(false);
  const [selectedProduct, setSelectedProduct] = createSignal(0);

  createEffect(() => {
    setMounted(true);
  });

  return (
    <div class="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
      <Motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: mounted() ? 1 : 0 }}
        transition={{ duration: 0.5 }}
        class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-24"
      >
        <div class="text-center mb-16">
          <Motion.h1
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            class="text-4xl sm:text-5xl md:text-6xl font-extrabold mb-4 text-gray-900"
          >
            Transforme Sua Gestão Empresarial
          </Motion.h1>
          <Motion.p
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.4 }}
            class="text-xl text-gray-600 max-w-3xl mx-auto mb-8"
          >
            Descubra ferramentas inteligentes e recursos completos para
            impulsionar o sucesso do seu negócio.
          </Motion.p>
          <Motion.button
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.6 }}
            onclick={() =>navigate("/register")}
            class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Comece Agora
            <ArrowRight class="ml-2 -mr-1 h-5 w-5" aria-hidden="true" />
          </Motion.button>
        </div>

        <div class="mb-16">
          <Motion.h2
            initial={{ y: -20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.5, delay: 0.8 }}
            class="text-3xl sm:text-4xl font-bold text-center mb-8 text-gray-800"
          >
            Nossas Soluções
          </Motion.h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <For each={products}>
              {(product, index) => (
                <Motion.div
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ duration: 0.5, delay: 0.2 * index() }}
                  class={`bg-white rounded-xl shadow-lg overflow-hidden cursor-pointer transition-all duration-300 hover:shadow-xl `}
                  onClick={() => setSelectedProduct(index())}
                >
                  <div class={`p-6 ${product.bgColor}`}>
                    <div class="flex items-center justify-between mb-4">
                      <h3 class="text-xl font-semibold text-gray-900">
                        {product.title}
                      </h3>
                      <product.icon size={32} class={product.color} />
                    </div>
                    <p class="text-gray-600">{product.description}</p>
                  </div>
                  <div class="p-6 bg-gray-50">
                    <h4 class="font-semibold text-gray-900 mb-2">Destaques:</h4>
                    <ul class="space-y-2">
                      <li class="flex items-start">
                        <Zap
                          class="text-yellow-500 mr-2 mt-1 flex-shrink-0"
                          size={18}
                        />
                        <span class="text-sm text-gray-600">
                          {product.performance}
                        </span>
                      </li>
                      <li class="flex items-start">
                        <Shield
                          class="text-green-500 mr-2 mt-1 flex-shrink-0"
                          size={18}
                        />
                        <span class="text-sm text-gray-600">
                          {product.security}
                        </span>
                      </li>
                    </ul>
                  </div>
                </Motion.div>
              )}
            </For>
          </div>
        </div>

        <Motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          class="bg-white rounded-xl shadow-xl overflow-hidden mb-16"
        >
          <div class="p-6 sm:p-10">
            <h2 class="text-3xl font-bold mb-4 text-gray-900">
              {products[selectedProduct()].title}
            </h2>
            <p class="text-gray-600 mb-6">
              {products[selectedProduct()].details}
            </p>
            <div class="grid sm:grid-cols-2 gap-6">
              <div class="flex items-start">
                <Zap
                  class="text-yellow-500 mr-3 mt-1 flex-shrink-0"
                  size={24}
                />
                <div>
                  <h3 class="text-lg font-semibold mb-2 text-gray-900">
                    Performance
                  </h3>
                  <p class="text-gray-600">
                    {products[selectedProduct()].performance}
                  </p>
                </div>
              </div>
              <div class="flex items-start">
                <Shield
                  class="text-green-500 mr-3 mt-1 flex-shrink-0"
                  size={24}
                />
                <div>
                  <h3 class="text-lg font-semibold mb-2 text-gray-900">
                    Segurança
                  </h3>
                  <p class="text-gray-600">
                    {products[selectedProduct()].security}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </Motion.div>

        <div class="text-center">
          <Motion.button
            animate={{ opacity: 1, y: 0 }}
            initial={{ opacity: 0, y: 20 }}
            onclick={() =>navigate("/register")}
            transition={{ duration: 0.5, delay: 0.2 }}
            class="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
            Experimente Gratuitamente
            <ArrowRight class="ml-2 -mr-1 h-5 w-5" aria-hidden="true" />
          </Motion.button>
        </div>
      </Motion.div>
    </div>
  );
};

export default Welcome;
