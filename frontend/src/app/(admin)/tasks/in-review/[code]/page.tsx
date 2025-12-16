import InReviewForm from "@/components/tasks/InReviewForm";

interface InReviewTaskPageProps {
  params: {
    code: string;
  };
}

export default function InReviewTaskPage({ params }: InReviewTaskPageProps) {
  return <InReviewForm code={params.code} />;
}
