export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen w-full bg-[#d8d8ec] relative">
      {/* Fixed Background Layer */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        {/* Wavy Pattern Background */}
        <div
          className="absolute inset-0 opacity-40"
          style={{
            backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 1200 800' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill='%237678ed' fill-opacity='0.12' d='M0,192L48,176C96,160,192,128,288,144C384,160,480,224,576,245.3C672,267,768,245,864,213.3C960,181,1056,139,1152,144C1248,149,1344,203,1392,229.3L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z'%3E%3C/path%3E%3C/svg%3E")`,
            backgroundSize: 'cover',
            backgroundPosition: 'bottom',
          }}
        />

        {/* Floating Orb 1 - Top Left (Purple) */}
        <div
          className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full pointer-events-none animate-float"
          style={{
            background: 'radial-gradient(circle, rgba(118,120,237,0.25) 0%, rgba(118,120,237,0) 70%)',
          }}
        />

        {/* Floating Orb 2 - Bottom Right (Coral) */}
        <div
          className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] rounded-full pointer-events-none animate-float-reverse"
          style={{
            background: 'radial-gradient(circle, rgba(255,122,85,0.2) 0%, rgba(255,122,85,0) 70%)',
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
