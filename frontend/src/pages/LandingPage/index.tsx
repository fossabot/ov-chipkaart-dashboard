import React from 'react';
import './index.scss';
import Grid from '@material-ui/core/Grid';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import ROUTE_NAMES from '../../constants/routes';
import { Link } from 'react-router-dom';
import Logo from '../../components/Logo';
import Typography from '@material-ui/core/Typography';
import { useTranslation } from 'react-i18next';
import TransKeys from '../../i18n/keys';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            height: '100vh',
            width: '100vw',
            overflow: 'hidden',
        },

        info: {
            height: '100vh',
            padding: theme.spacing(2),
            backgroundColor: theme.palette.primary.dark,
            display: 'flex',
            alignItems: 'center',
        },

        auth: {
            height: '100vh',
            display: 'flex',
            alignItems: 'center',
            padding: theme.spacing(3),
            backgroundColor: theme.palette.primary.contrastText,
        },

        card: {
            color: theme.palette.primary.contrastText,
        },

        logo: {
            position: 'absolute',
        },
    }),
);

export default function LandingPage() {
    const classes = useStyles();
    const { t } = useTranslation();

    return (
        <Grid container className={classes.container}>
            <Grid item xs={8}>
                <Link className={classes.logo} to={ROUTE_NAMES.DASHBOARD}>
                    <Logo />
                </Link>
                <Box className={classes.info}>
                    <Box className={classes.card} width="70%" marginTop="300">
                        <Typography variant="h2">
                            {t(TransKeys.LANDING_PAGE.TITLE)}
                        </Typography>
                        <Typography variant="subtitle1" gutterBottom>
                            {t(TransKeys.LANDING_PAGE.SUB_TITLE)}
                        </Typography>
                    </Box>
                </Box>
            </Grid>
            <Grid item xs={4}>
                <Box className={classes.auth}>Second</Box>
            </Grid>
        </Grid>
    );
}
