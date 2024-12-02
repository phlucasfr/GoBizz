import { BarChart3 } from "lucide-solid";

type QuickStatsCardProps = {
    title: string;
    value: string;
    change: string;
  };

const QuickStatsCard = (props: QuickStatsCardProps) => (
    <div class="bg-white p-6 rounded-lg shadow-md">
      <div class="flex items-center justify-between space-y-0 pb-2">
        <h3 class="text-sm font-medium">{props.title}</h3>
        <BarChart3 class="h-4 w-4 text-gray-400" />
      </div>
      <div>
        <div class="text-2xl font-bold">{props.value}</div>
        <p class={`text-xs ${props.change.startsWith('+') ? 'text-green-500' : 'text-red-500'}`}>
          {props.change}
        </p>
      </div>
    </div>
  );

export default QuickStatsCard;