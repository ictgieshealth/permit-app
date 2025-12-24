import ShowTask from "@/components/tasks/ShowTask";

interface ShowTaskPageProps {
  params: Promise<{
    code: string;
  }>;
}

export default async function ShowTaskPage({ params }: ShowTaskPageProps) {
  const { code } = await params;
  return <ShowTask code={code} />;
}
