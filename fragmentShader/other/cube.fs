#version 330 core
in vec3 uv;
uniform samplerCube cubemap;
uniform mat4 transform;
out vec4 FragColor;
void main()
{
    FragColor = texture(cubemap, uv);
    FragColor += transform * uv.x * uv.y * uv.z * vec4(1.0,0.0,0.0,0.0);
}