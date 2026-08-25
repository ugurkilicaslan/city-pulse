/**
 * City Pulse — Starfield Canvas Animation
 * Call: initStarfield('canvas-id', options)
 */
(function() {

    function initStarfield(canvasId, opts) {
        opts = opts || {};
        const canvas  = document.getElementById(canvasId);
        if (!canvas) return;
        const ctx     = canvas.getContext('2d');
        const COUNT   = opts.count  || 180;
        const SPEED   = opts.speed  || 0.25;
        const OPACITY = opts.opacity || 0.7;

        function resize() {
            canvas.width  = window.innerWidth;
            canvas.height = window.innerHeight;
        }
        resize();
        window.addEventListener('resize', resize);

        /* Stars */
        const stars = Array.from({ length: COUNT }, () => ({
            x:  Math.random(),
            y:  Math.random(),
            r:  Math.random() * 1.4 + 0.2,
            vx: (Math.random() - 0.5) * SPEED * 0.4,
            vy: (Math.random() - 0.5) * SPEED * 0.4,
            op: Math.random() * 0.5 + 0.15,
            ph: Math.random() * Math.PI * 2,
            ps: Math.random() * 0.015 + 0.004,
        }));

        /* Nebula blobs */
        const blobs = [
            { xr: 0.15, yr: 0.25, r: 380, color: [109, 40, 217], a: 0.07 },
            { xr: 0.82, yr: 0.70, r: 300, color: [76,  29, 149], a: 0.06 },
            { xr: 0.50, yr: 0.05, r: 250, color: [255,197,  61], a: 0.03 },
            { xr: 0.70, yr: 0.40, r: 200, color: [45,  212,191], a: 0.03 },
        ];

        /* Shooting stars */
        const shooters = [];
        function spawnShooter() {
            shooters.push({
                x: Math.random() * canvas.width,
                y: Math.random() * canvas.height * 0.5,
                len: Math.random() * 120 + 60,
                vx: Math.random() * 6 + 4,
                vy: Math.random() * 3 + 1,
                life: 1,
                decay: Math.random() * 0.02 + 0.015,
            });
        }
        setInterval(spawnShooter, 2800);

        let t = 0;
        function draw() {
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            /* Nebula blobs */
            blobs.forEach(b => {
                const pulse = Math.sin(t * 0.004 + b.xr * 5) * 0.2 + 0.8;
                const gx = b.xr * canvas.width;
                const gy = b.yr * canvas.height;
                const g = ctx.createRadialGradient(gx, gy, 0, gx, gy, b.r * pulse);
                g.addColorStop(0, `rgba(${b.color[0]},${b.color[1]},${b.color[2]},${b.a})`);
                g.addColorStop(1, 'transparent');
                ctx.fillStyle = g;
                ctx.fillRect(0, 0, canvas.width, canvas.height);
            });

            /* Stars */
            stars.forEach(s => {
                const pulse = Math.sin(t * s.ps + s.ph) * 0.35 + 0.65;
                const alpha = s.op * pulse * OPACITY;
                ctx.beginPath();
                ctx.arc(s.x * canvas.width, s.y * canvas.height, s.r * pulse, 0, Math.PI * 2);
                ctx.fillStyle = `rgba(255,255,255,${alpha})`;
                ctx.fill();

                s.x += s.vx / canvas.width;
                s.y += s.vy / canvas.height;
                if (s.x < 0) s.x = 1;
                if (s.x > 1) s.x = 0;
                if (s.y < 0) s.y = 1;
                if (s.y > 1) s.y = 0;
            });

            /* Shooting stars */
            for (let i = shooters.length - 1; i >= 0; i--) {
                const sh = shooters[i];
                ctx.beginPath();
                ctx.moveTo(sh.x, sh.y);
                ctx.lineTo(sh.x - sh.len * sh.vx / 8, sh.y - sh.len * sh.vy / 8);
                const grad = ctx.createLinearGradient(
                    sh.x, sh.y,
                    sh.x - sh.len * sh.vx / 8, sh.y - sh.len * sh.vy / 8
                );
                grad.addColorStop(0, `rgba(255,255,255,${sh.life * 0.9})`);
                grad.addColorStop(1, 'transparent');
                ctx.strokeStyle = grad;
                ctx.lineWidth = 1.5;
                ctx.stroke();
                sh.x += sh.vx;
                sh.y += sh.vy;
                sh.life -= sh.decay;
                if (sh.life <= 0) shooters.splice(i, 1);
            }

            t++;
            requestAnimationFrame(draw);
        }
        draw();
    }

    window.initStarfield = initStarfield;

})();
