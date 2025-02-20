mod handlers;
mod utils;
use poise::serenity_prelude as serenity;

#[shuttle_runtime::main]
async fn main(
    #[shuttle_shared_db::Postgres] pool: sqlx::PgPool,
    #[shuttle_runtime::Secrets] secret_store: shuttle_runtime::SecretStore,
) -> shuttle_serenity::ShuttleSerenity {
    sqlx::migrate!()
        .run(&pool)
        .await
        .expect("Failed to run migrations on database");

    // Get the discord token set in `Secrets.toml`
    let discord_token = secret_store
        .get("DISCORD_TOKEN")
        .expect("DISCORD_TOKEN not set");

    let framework = poise::Framework::builder()
        .options(poise::FrameworkOptions {
            commands: vec![handlers::admins::admins(), handlers::counting::counting()],
            event_handler: |context, event, framework, data| {
                Box::pin(handlers::event::event_handler(
                    context, event, framework, data,
                ))
            },
            ..Default::default()
        })
        .setup(|context, _ready, framework| {
            Box::pin(async move {
                poise::builtins::register_globally(context, &framework.options().commands).await?;

                Ok(utils::ServerData {
                    cooldowns: Default::default(),
                    pool,
                })
            })
        })
        .build();

    let client = serenity::ClientBuilder::new(
        discord_token,
        serenity::GatewayIntents::non_privileged() | serenity::GatewayIntents::MESSAGE_CONTENT,
    )
    .framework(framework)
    .await
    .map_err(shuttle_runtime::CustomError::new)?;

    Ok(client.into())
}
