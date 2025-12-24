import TaskForm from "@/components/tasks/TaskForm";

interface EditTaskPageProps {
  params: {
    code: string;
  };
}

export default async function EditTaskPage({ params }: EditTaskPageProps) {
  const { code } = await params;
  return <TaskForm code={code} isEdit={true} />;
}
