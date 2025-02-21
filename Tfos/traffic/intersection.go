package traffic

import (
    "context"
    "encoding/json"
    "time"
)

// AdjustIntersection modifies traffic light timing based on alerts
func AdjustIntersection(id string, alert map[string]interface{}) error {
    // Convert alert to JSON
    alertJSON, err := json.Marshal(alert)
    if err != nil {
        return err
    }

    // Store alert in Redis
    err = rdb.HSet(context.Background(), "intersection:"+id,
        "alert_active", "true",
        "alert_data", string(alertJSON),
        "alert_timestamp", time.Now().Format(time.RFC3339),
    ).Err()
    if err != nil {
        return err
    }

    // Trigger immediate optimization for this intersection
    cfg := getIntersectionConfig(id)
    newPhase := calculateOptimalPhase(cfg)
    applyPhase(cfg, newPhase)

    return nil
}