import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { toast } from "sonner";

export function RegisterForm({ className, onSwitchToLogin, onRegister, ...props }) {

  const handleSubmit = (e) => {
    e.preventDefault();
    if (e.target.username.value === "" || e.target.password.value === "") {
      toast.error("Username and password are required.", {
        position: "top-right",
        style: {
          background: "white",
          color: "red",
        },
      });
      return;
    }
    onRegister({
      username: e.target.username.value,
      password: e.target.password.value,
    });
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <CardTitle>Register</CardTitle>
          <CardDescription>
            Enter your username and password below to register
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="username">User Name</Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="username"
                  required
                />
              </div>
              <div className="grid gap-3">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                </div>
                <Input id="password" type="password" required />
              </div>
              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full">
                  Register
                </Button>

                <div className="flex flex-col gap-3 text-center hover:underline">
                  <a href="#" onClick={() => onSwitchToLogin()}>Already have an account? Login</a>
                </div>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

export default RegisterForm;
