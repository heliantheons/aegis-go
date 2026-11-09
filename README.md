<p align="center">
  <img src="./assets/brand/hero-ice.png" width="256" alt="Aegis Go logo" />
</p>

<h1 align="center">Aegis Go</h1>

Aegis Go 是一组跟 Aegis 打交道的 Go 库。没有单一入口，而是按用途拆成四个各自独立发版的 module——`application`（应用侧注册与令牌管理）、`guard`（认证中间件与鉴权 requirement）、`service`（令牌的提取、校验、解密）、`utilities`（密钥、PASETO/PASERK、issuer 等基础能力）。需要哪一层就只依赖哪一层，不必把整套都拉进来。

Aegis Go is a set of Go libraries for working with Aegis. Rather than a single entry point, it's split into four independently versioned modules — `application`, `guard`, `service`, and `utilities` — so you only depend on the layer you actually use.

## 服务访问令牌

SAT 使用签名内的 `token_type: "sat"` 标识令牌类型。它与 Challenge 的业务 `typ` 字段相互独立；解析器优先使用显式令牌类型，再使用其他令牌既有的判定规则。SAT 的明文 `sub` 不会进入 UAT 的身份解密流程。

签发方和消费方应一起升级 `utilities`，再重新签发服务令牌。旧的无类型 SAT 无法与 UAT 或 CT 可靠区分，不应用来调用要求 SAT 的接口。

Service access tokens carry a signed `token_type: "sat"` claim, independently of a challenge's business `typ`. Upgrade the issuer and consumer utilities together and issue fresh service tokens; the SDK does not guess a service identity from unmarked user or client claims.