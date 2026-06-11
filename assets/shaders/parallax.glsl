// Parallax Occlusion Mapping surface shader for R3D.
// Offsets the UVs along the tangent-space view direction using a height map,
// then samples the material (albedo/normal/ORM) at the displaced UVs.
#define R3D_NO_AUTO_FETCH

uniform sampler2D u_height;   // grayscale height (1 = top surface)
uniform vec3  u_view_pos;     // camera world position
uniform float u_depth_scale;  // displacement strength (e.g. 0.07)
uniform float u_tile;         // UV tiling factor
uniform int   u_enable;       // 0 = plain normal mapping, 1 = parallax

vec2 parallax(vec2 uv, vec3 viewT) {
    const float maxL = 48.0, minL = 8.0;
    float numLayers = mix(maxL, minL, clamp(abs(viewT.z), 0.0, 1.0));
    float layerDepth = 1.0 / numLayers;
    vec2 P = (viewT.xy / viewT.z) * u_depth_scale;
    vec2 dUV = P / numLayers;

    float curLayer = 0.0;
    vec2  cuv = uv;
    float curH = 1.0 - texture(u_height, cuv).r;
    for (int i = 0; i < 64; i++) {
        if (curLayer >= curH) break;
        cuv -= dUV;
        curH = 1.0 - texture(u_height, cuv).r;
        curLayer += layerDepth;
    }
    // occlusion interpolation between the last two layers
    vec2  prevUV = cuv + dUV;
    float after  = curH - curLayer;
    float before = (1.0 - texture(u_height, prevUV).r) - (curLayer - layerDepth);
    float w = after / (after - before);
    return mix(cuv, prevUV, clamp(w, 0.0, 1.0));
}

void fragment() {
    vec2 uv = TEXCOORD * u_tile;
    if (u_enable != 0) {
        mat3 TBN = mat3(normalize(TANGENT), normalize(BITANGENT), normalize(NORMAL));
        vec3 viewT = normalize(transpose(TBN) * (u_view_pos - POSITION));
        if (viewT.z > 0.05) uv = parallax(uv, viewT);
    }
    FetchMaterial(uv);
}
