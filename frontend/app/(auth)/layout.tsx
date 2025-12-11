export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen w-full bg-slate-950 relative">
      {/* Fixed Background Layer */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        {/* Animated Gradient Mesh Background */}
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(120,119,198,0.3),rgba(255,255,255,0))]" />

        {/* Floating Orb 1 - Top Left (Indigo) */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(99,102,241,0.3) 0%, rgba(99,102,241,0) 70%)',
          }}
        />

        {/* Floating Orb 2 - Bottom Right (Purple) */}
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(168,85,247,0.25) 0%, rgba(168,85,247,0) 70%)',
          }}
        />
      </div>

      {/* Content */}
      <div className="relative z-10 w-full h-full animate-enter">
        {children}
      </div>
    </div>
  );
}
