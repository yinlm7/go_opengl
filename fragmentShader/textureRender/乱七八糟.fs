#version 330 core
out vec4 FragColor;

uniform sampler2D texture0;
uniform float time;
in vec2 uv;


mat2 rotate2d(float a) {
    return mat2(cos(a), -sin(a), sin(a), cos(a));
}

void main() {
    vec2 fragCoord = gl_FragCoord.xy;
    vec2 iResolution = fragCoord/uv;

    float size = iResolution.x / (25.0 + pow(sin(time/5.5),9.)*10.);
    float angle = time / 10.0;
    float kx = sin(time/1.)/(iResolution.x*1.3*(1.-pow(sin(time/15.),5.)));
    float ky = pow(sin(time/1.4),3.)/(iResolution.x*1.3*(1.-pow(sin(time/22.),7.)));
    fragCoord = fragCoord.xy - iResolution.xy / 2.0;
    float stretch_x = fragCoord.x - iResolution.x * 1.5 / 2.0 * sin(time/11.);
    float stretch_y = fragCoord.y - iResolution.y * 1.5 / 2.0 * sin(time/19.);
    vec2 p = fragCoord.xy * rotate2d(angle + kx*stretch_x - ky*stretch_y);
    vec2 pmod = vec2(mod(p.x, size * 2.0), mod(p.y, size * 2.0));
    float k = max(pmod.x, pmod.y)/size/2.;

    FragColor = texture2D(texture0, p/iResolution.xy - vec2(0.5,0.5));
}