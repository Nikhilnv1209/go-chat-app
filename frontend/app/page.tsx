import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center p-8 bg-background">
      <Card className="w-full max-w-md text-center border-none shadow-none">
        <CardHeader>
          <CardTitle className="text-4xl font-bold tracking-tight text-primary">Go Chat</CardTitle>
          <CardDescription className="text-lg mt-2">
            Real-time messaging, redefined.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col gap-4 mt-8">
          <div className="flex gap-4 justify-center">
            <Button asChild size="lg" className="w-full">
              <Link href="/login">Login</Link>
            </Button>
            <Button asChild variant="outline" size="lg" className="w-full">
              <Link href="/register">Register</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
