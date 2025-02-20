const COOLDOWN_DURATION: std::time::Duration = std::time::Duration::from_secs(5);

pub fn add(data: &crate::utils::ServerData, author_id: i64) {
    let mut cooldowns = data.cooldowns.lock().unwrap();
    cooldowns.insert(author_id, std::time::Instant::now());
}

pub fn check_if_on(data: &crate::utils::ServerData, author_id: i64) -> bool {
    let mut cooldowns = data.cooldowns.lock().unwrap();
    if let Some(last_used) = cooldowns.get(&author_id) {
        if last_used.elapsed() < COOLDOWN_DURATION {
            return true;
        } else {
            cooldowns.remove(&author_id);
        }
    }

    false
}
