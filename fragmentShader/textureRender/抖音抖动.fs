#version 330 core
out vec4 FragColor;

uniform sampler2D texture;
in vec2 uv;
uniform float time;

const float duration = 0.7;
const float maxScale = 1.1;
const float offset   = 0.02;

void main() {

    float progress = mod(time, duration) / duration;
    vec2 offsetCoords = vec2(offset, offset) * progress;

    float scale = 1.0 + (maxScale - 1.0) * progress;
    vec2 ScaleTextureCoords = vec2(0.5, 0.5) + (uv - vec2(0.5, 0.5)) / scale;

    vec4 maskR = texture2D(texture, ScaleTextureCoords + offsetCoords);
    vec4 maskB = texture2D(texture, ScaleTextureCoords + offsetCoords);
    vec4 mask  = texture2D(texture, ScaleTextureCoords);

    FragColor = vec4(maskR.r, mask.g, maskB.b, mask.a);
}