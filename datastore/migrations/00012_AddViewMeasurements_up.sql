-- bactdb
-- Matthew R Dillon

CREATE OR REPLACE VIEW v_measurements AS
SELECT s.strain_name,
    o.observation_name,
    som.measurement_table,
    tm.text_measurement_name,
    nm.measurement_value,
    nm.confidence_interval,
    nm.unit_type_id
FROM strainsobsmeasurements som
INNER JOIN strainsobservations so
    ON som.strainsobservations_id = so.id
INNER JOIN strains s
    ON so.strain_id = s.id
INNER JOIN observations o
    ON so.observations_id = o.id
LEFT OUTER JOIN text_measurements tm
    ON som.measurement_id = tm.id
    AND som.measurement_table = 'text'
LEFT OUTER JOIN numerical_measurements nm
    ON som.measurement_id = nm.id
    AND som.measurement_table = 'num'
ORDER BY measurement_table, o.observation_name ASC

